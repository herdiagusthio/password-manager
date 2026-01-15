package domain

import (
	"context"
	"time"
)

type Secret struct {
	ID                string    `json:"id"`
	UserID            string    `json:"user_id"`
	Title             string    `json:"title"`
	Username          string    `json:"username"`
	EncryptedPassword string    `json:"-"` // Never expose directly in JSON without decryption
	Password          string    `json:"password,omitempty"` // Decrypted password, only populated when needed
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	Version           int       `json:"version"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type SecretRepository interface {
	Create(ctx context.Context, secret *Secret) error
	GetByID(ctx context.Context, id string) (*Secret, error)
	ListByUserID(ctx context.Context, userID string) ([]*Secret, error)
	Update(ctx context.Context, secret *Secret) error
	Delete(ctx context.Context, id string) error
}

type SecretUsecase interface {
	CreateSecret(ctx context.Context, secret *Secret) error
	GetSecret(ctx context.Context, id string, userID string) (*Secret, error)
	ListSecrets(ctx context.Context, userID string) ([]*Secret, error)
	UpdateSecret(ctx context.Context, secret *Secret) error
	DeleteSecret(ctx context.Context, id string, userID string) error
}
