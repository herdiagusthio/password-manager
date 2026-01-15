package usecase

import (
	"context"
	"fmt"

	"github.com/herdiagusthio/password-manager/config"
	"github.com/herdiagusthio/password-manager/internal/domain"
	"github.com/herdiagusthio/password-manager/pkg/crypto"
)

type secretUsecase struct {
	repo domain.SecretRepository
	cfg  *config.Config
}

func NewSecretUsecase(repo domain.SecretRepository, cfg *config.Config) domain.SecretUsecase {
	return &secretUsecase{
		repo: repo,
		cfg:  cfg,
	}
}

func (u *secretUsecase) CreateSecret(ctx context.Context, secret *domain.Secret) error {
	// Encrypt the password before saving
	// Note: We use the system Master Key from config for simplicity in this version.
	// Production apps might use a per-user key derived from KDF or KMS.
	encrypted, err := crypto.Encrypt(secret.Password, u.cfg.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}
	secret.EncryptedPassword = encrypted

	// Clear plain password from struct to avoid accidental leak later
	secret.Password = "" 

	return u.repo.Create(ctx, secret)
}

func (u *secretUsecase) GetSecret(ctx context.Context, id string, userID string) (*domain.Secret, error) {
	secret, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, nil // Not found
	}

	// Authorization check
	if secret.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to secret")
	}

	// Decrypt
	decrypted, err := crypto.Decrypt(secret.EncryptedPassword, u.cfg.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt password: %w", err)
	}
	secret.Password = decrypted

	return secret, nil
}

func (u *secretUsecase) ListSecrets(ctx context.Context, userID string) ([]*domain.Secret, error) {
	// We list secrets but do NOT return the decrypted passwords in the list view for security/performance
	return u.repo.ListByUserID(ctx, userID)
}

func (u *secretUsecase) UpdateSecret(ctx context.Context, secret *domain.Secret) error {
	// Check existance and ownership first
	existing, err := u.repo.GetByID(ctx, secret.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("secret not found")
	}
	if existing.UserID != secret.UserID {
		return fmt.Errorf("unauthorized update")
	}

	// If a new password is provided, encrypt it. Otherwise keep existing.
	if secret.Password != "" {
		encrypted, err := crypto.Encrypt(secret.Password, u.cfg.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt password: %w", err)
		}
		secret.EncryptedPassword = encrypted
		secret.Password = ""
	} else {
		secret.EncryptedPassword = existing.EncryptedPassword
	}

	return u.repo.Update(ctx, secret)
}

func (u *secretUsecase) DeleteSecret(ctx context.Context, id string, userID string) error {
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return nil // Already gone
	}
	if existing.UserID != userID {
		return fmt.Errorf("unauthorized delete")
	}

	return u.repo.Delete(ctx, id)
}
