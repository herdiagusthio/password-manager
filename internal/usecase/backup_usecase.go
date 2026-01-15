package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/herdiagusthio/password-manager/config"
	"github.com/herdiagusthio/password-manager/internal/domain"
	"github.com/herdiagusthio/password-manager/pkg/crypto"
)

type backupUsecase struct {
	secretRepo domain.SecretRepository
	cfg        *config.Config
}

func NewBackupUsecase(secretRepo domain.SecretRepository, cfg *config.Config) domain.BackupUsecase {
	return &backupUsecase{
		secretRepo: secretRepo,
		cfg:        cfg,
	}
}

func (u *backupUsecase) ExportSecrets(ctx context.Context, userID string) ([]byte, error) {
	// 1. Fetch all secrets
	secrets, err := u.secretRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}

	// 2. Prepare Backup struct
	backup := domain.Backup{
		Version:   "1.0",
		CreatedAt: time.Now().Format(time.RFC3339),
		Secrets:   secrets,
	}

	// 3. Serialize to JSON
	jsonData, err := json.Marshal(backup)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal backup: %w", err)
	}

	// 4. Encrypt the entire JSON blob with Master Key
	encryptedString, err := crypto.Encrypt(string(jsonData), u.cfg.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt backup: %w", err)
	}

	return []byte(encryptedString), nil
}

func (u *backupUsecase) ImportSecrets(ctx context.Context, userID string, backupData []byte) error {
	// 1. Decrypt
	decryptedJSON, err := crypto.Decrypt(string(backupData), u.cfg.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt backup: %w", err)
	}

	// 2. Unmarshal
	var backup domain.Backup
	if err := json.Unmarshal([]byte(decryptedJSON), &backup); err != nil {
		return fmt.Errorf("failed to unmarshal backup json: %w", err)
	}

	// 3. Restore
	// Strategy: Iterate and create/update.
	// We will overwrite existing secrets with same Title? Or just add new ones?
	// The robust way for "Restore" is usually: Upsert based on ID if present, or create new if not.
	// Since IDs in backup are UUIDs, if they match, we update.

	for _, s := range backup.Secrets {
		// Enforce UserID to be the current user (prevent restoring secrets to wrong user if backup file is shared/hacked)
		s.UserID = userID
		
		// Check if exists
		existing, err := u.secretRepo.GetByID(ctx, s.ID)
		if err != nil {
			return fmt.Errorf("error checking secret existence: %w", err)
		}

		if existing != nil {
			if err := u.secretRepo.Update(ctx, s); err != nil {
				return fmt.Errorf("failed to update secret %s: %w", s.ID, err)
			}
		} else {
			if err := u.secretRepo.Create(ctx, s); err != nil {
				return fmt.Errorf("failed to create secret %s: %w", s.ID, err)
			}
		}
	}

	return nil
}
