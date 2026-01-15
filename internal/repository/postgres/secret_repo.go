package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/herdiagusthio/password-manager/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type secretRepo struct {
	db *pgxpool.Pool
}

func NewSecretRepository(db *pgxpool.Pool) domain.SecretRepository {
	return &secretRepo{
		db: db,
	}
}

func (r *secretRepo) Create(ctx context.Context, secret *domain.Secret) error {
	query := `
		INSERT INTO secrets (user_id, title, username, encrypted_password, metadata, version)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	row := r.db.QueryRow(ctx, query,
		secret.UserID,
		secret.Title,
		secret.Username,
		secret.EncryptedPassword,
		secret.Metadata,
		secret.Version,
	)

	err := row.Scan(&secret.ID, &secret.CreatedAt, &secret.UpdatedAt)
	if err != nil {
		return fmt.Errorf("secretRepo.Create: %w", err)
	}
	return nil
}

func (r *secretRepo) GetByID(ctx context.Context, id string) (*domain.Secret, error) {
	query := `
		SELECT id, user_id, title, username, encrypted_password, metadata, version, created_at, updated_at
		FROM secrets
		WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)

	var s domain.Secret
	err := row.Scan(
		&s.ID, &s.UserID, &s.Title, &s.Username, &s.EncryptedPassword, &s.Metadata, &s.Version, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("secretRepo.GetByID: %w", err)
	}
	return &s, nil
}

func (r *secretRepo) ListByUserID(ctx context.Context, userID string) ([]*domain.Secret, error) {
	query := `
		SELECT id, user_id, title, username, encrypted_password, metadata, version, created_at, updated_at
		FROM secrets
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("secretRepo.ListByUserID query: %w", err)
	}
	defer rows.Close()

	var secrets []*domain.Secret
	for rows.Next() {
		var s domain.Secret
		err := rows.Scan(
			&s.ID, &s.UserID, &s.Title, &s.Username, &s.EncryptedPassword, &s.Metadata, &s.Version, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("secretRepo.ListByUserID scan: %w", err)
		}
		secrets = append(secrets, &s)
	}
	return secrets, nil
}

func (r *secretRepo) Update(ctx context.Context, secret *domain.Secret) error {
	query := `
		UPDATE secrets
		SET title = $1, username = $2, encrypted_password = $3, metadata = $4, version = version + 1, updated_at = NOW()
		WHERE id = $5
		RETURNING version, updated_at
	`
	row := r.db.QueryRow(ctx, query,
		secret.Title,
		secret.Username,
		secret.EncryptedPassword,
		secret.Metadata, // Metadata is interface{}, pgx handles JSONB mapping
		secret.ID,
	)

	err := row.Scan(&secret.Version, &secret.UpdatedAt)
	if err != nil {
		return fmt.Errorf("secretRepo.Update: %w", err)
	}
	return nil
}

func (r *secretRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM secrets WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("secretRepo.Delete: %w", err)
	}
	return nil
}
