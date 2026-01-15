package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/herdiagusthio/password-manager/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) domain.AuthRepository {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, created_at, updated_at FROM users WHERE email = $1`
	row := r.db.QueryRow(ctx, query, email)

	var user domain.User
	err := row.Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Not found is not an error in this context, just nil
		}
		return nil, fmt.Errorf("userRepo.GetByEmail: %w", err)
	}
	return &user, nil
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (email) VALUES ($1) RETURNING id, created_at, updated_at`
	row := r.db.QueryRow(ctx, query, user.Email)
	
	err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("userRepo.Create: %w", err)
	}
	return nil
}
