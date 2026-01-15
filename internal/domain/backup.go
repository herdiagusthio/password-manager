package domain

import (
	"context"
)

// Backup represents the structure of the exported file
type Backup struct {
	Version   string    `json:"version"`
	CreatedAt string    `json:"created_at"`
	Secrets   []*Secret `json:"secrets"`
}

type BackupUsecase interface {
	// ExportSecrets returns an encrypted byte slice of the backup
	ExportSecrets(ctx context.Context, userID string) ([]byte, error)
	
	// ImportSecrets takes an encrypted byte slice and restores secrets
	ImportSecrets(ctx context.Context, userID string, backupData []byte) error
}
