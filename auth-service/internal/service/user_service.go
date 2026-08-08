package service

import (
	"bookshelf/auth-service/internal/domain"
	"bookshelf/auth-service/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrUsernameExists     = errors.New("username already exists")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidUsername    = errors.New("invalid username")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrUserNotFound       = errors.New("user not found")
)

type UserService struct {
	repo      *repository.UserRepository
	jwtSecret string
}

func NewUserService(userRepo *repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		repo:      userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *UserService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error) {
	if err := s.validateRegister(ctx, req); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	if err = s.repo.Create(ctx, &user); err != nil {
		return nil, err
	}

	return s.createAuthResponse(&user)
}

func (s *UserService) validateRegister(ctx context.Context, req domain.RegisterRequest) error {
	if len(req.Email) == 0 {
		return ErrInvalidEmail
	}

	if len(req.Username) < 3 {
		return ErrInvalidUsername
	}

	if len(req.Password) < 8 {
		return ErrInvalidPassword
	}

	if s.repo.EmailExists(ctx, req.Email) {
		return ErrUserExists
	}

	if s.repo.UsernameExists(ctx, req.Username) {
		return ErrUsernameExists
	}

	return nil
}

func (s *UserService) generateAccessToken(userID string) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(now),
	})

	return token.SignedString([]byte(s.jwtSecret))
}

func (s *UserService) createAuthResponse(user *domain.User) (*domain.AuthResponse, error) {
	token, err := s.generateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(time.Now().Add(time.Hour).Unix()),
		User:        user.ToPublic(),
	}, nil
}

func (s *UserService) ValidateToken(tokenString string) (string, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims.Subject, nil
}

func (s *UserService) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.createAuthResponse(user)
}

func (s *UserService) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, repository.ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID string, req domain.UpdateUserRequest) (*domain.User, error) {
	if len(req.Username) < 3 {
		return nil, ErrInvalidUsername
	}

	checkUser, err := s.repo.GetByUsername(ctx, req.Username)
	switch {
	case errors.Is(err, repository.ErrUserNotFound):
	case err != nil:
		return nil, err
	case checkUser.ID != userID:
		return nil, ErrUsernameExists
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.Username = req.Username

	if err = s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
