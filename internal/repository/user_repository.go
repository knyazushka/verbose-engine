// internal/repository/user_repository.go
package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knyazushka/verbose-engine/internal/domain"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		insert into users (email, username, password_hash, is_active)
		values ($1, $2, $3, $4)
		returning id, created_at, updated_at
	`

	return r.pool.QueryRow(ctx, query,
		user.Email, user.Username, user.PasswordHash, user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		select u.id, u.email, u.username, u.is_active, u.created_at, u.updated_at, u.deleted_at
		from users u
		where u.email = $1
	`

	var user domain.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, userId string) (*domain.User, error) {
	query := `
		select u.id, u.email, u.username, u.is_active, u.created_at, u.updated_at, u.deleted_at
		from users u
		where u.id = $1
	`

	var user domain.User
	err := r.pool.QueryRow(ctx, query, userId).Scan(
		&user.ID, &user.Email, &user.Username, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `select exists (select 1 from users where email = $1)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, email).Scan(&exists)
	return exists, err
}
