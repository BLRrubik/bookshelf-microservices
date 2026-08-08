package repository

import (
	"bookshelf/auth-service/internal/domain"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	createUserQuery = `
INSERT INTO users (id, username, email, password_hash)
VALUES ($1, $2, $3, $4)
ON CONFLICT DO NOTHING;
`
	getUserByIDQuery = `
SELECT id, username, email, password_hash, created_at, updated_at
FROM users
WHERE id = $1;
`
	getUserByUsernameQuery = `
SELECT id, username, email, password_hash, created_at, updated_at
FROM users
WHERE username = $1;
`
	getUserByEmailQuery = `
SELECT id, username, email, password_hash, created_at, updated_at
FROM users
WHERE email = $1;
`
	updateUserQuery = `
UPDATE users
SET username = $1, email = $2, password_hash = $3
WHERE id = $4;
`
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) Create(ctx context.Context, user *domain.User) error {
	user.ID = uuid.NewString()

	_, err := ur.db.ExecContext(ctx, createUserQuery, user.ID, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := ur.db.GetContext(ctx, &user, getUserByIDQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := ur.db.GetContext(ctx, &user, getUserByUsernameQuery, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := ur.db.GetContext(ctx, &user, getUserByEmailQuery, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) Update(ctx context.Context, user *domain.User) error {
	_, err := ur.db.ExecContext(ctx, updateUserQuery, user.ID, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		return err
	}

	return nil
}
