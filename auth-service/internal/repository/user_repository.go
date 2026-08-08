package repository

import (
	"bookshelf/auth-service/internal/domain"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
	existsUserByUsernameQuery = `
SELECT EXISTS(SELECT id FROM users WHERE username = $1);
`

	existsUserByEmailQuery = `
SELECT EXISTS(SELECT id FROM users WHERE email = $1);
`
	getUsersByIDsQuery = `
SELECT id, username, email, password_hash, created_at, updated_at
FROM users
WHERE id = ANY ($1);
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

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	user.ID = uuid.NewString()

	_, err := r.db.ExecContext(ctx, createUserQuery, user.ID, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.db.GetContext(ctx, &user, getUserByIDQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := r.db.GetContext(ctx, &user, getUserByUsernameQuery, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.GetContext(ctx, &user, getUserByEmailQuery, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(ctx, updateUserQuery, user.ID, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) UsernameExists(ctx context.Context, username string) bool {
	var exists bool
	err := r.db.GetContext(ctx, &exists, existsUserByUsernameQuery, username)
	if err != nil {
		return false
	}

	return exists
}

func (r *UserRepository) EmailExists(ctx context.Context, email string) bool {
	var exists bool
	err := r.db.GetContext(ctx, &exists, existsUserByEmailQuery, email)
	if err != nil {
		return false
	}

	return exists
}

func (r *UserRepository) GetByIDs(ctx context.Context, ids []string) ([]domain.User, error) {
	var users []domain.User

	if err := r.db.SelectContext(ctx, &users, getUsersByIDsQuery, pq.Array(ids)); err != nil {
		return nil, err
	}

	return users, nil
}
