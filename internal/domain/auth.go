package domain

import (
	"context"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuthRepository defines persistence methods for User
type AuthRepository interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
}

// AuthUsecase defines business logic for Authentication
type AuthUsecase interface {
	GetLoginURL(state string) string
	HandleCallback(ctx context.Context, code string) (*User, error)
	// Additional methods for Session management could go here
}
