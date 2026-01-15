package integration

import (
	"context"
	"testing"

	"github.com/herdiagusthio/password-manager/internal/domain"
	"github.com/herdiagusthio/password-manager/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecretRepo(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: database not initialized")
	}

	// Setup: Need a user first
	userRepo := postgres.NewUserRepository(testDB)
	secretRepo := postgres.NewSecretRepository(testDB)
	ctx := context.Background()

	user := &domain.User{Email: "secretagent@example.com"}
	require.NoError(t, userRepo.Create(ctx, user))

	t.Run("CreateAndGetSecret", func(t *testing.T) {
		secret := &domain.Secret{
			UserID:            user.ID,
			Title:             "Gmail",
			Username:          "agent",
			EncryptedPassword: "enc_password_123",
			Metadata:          map[string]interface{}{"url": "gmail.com"},
		}

		err := secretRepo.Create(ctx, secret)
		require.NoError(t, err)
		assert.NotEmpty(t, secret.ID)

		// Fetch
		found, err := secretRepo.GetByID(ctx, secret.ID)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, secret.Title, found.Title)
		assert.Equal(t, secret.Metadata["url"], found.Metadata["url"])
	})

	t.Run("ListSecrets", func(t *testing.T) {
		// Should see the one created above
		secrets, err := secretRepo.ListByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(secrets), 1)
	})

	t.Run("UpdateSecret", func(t *testing.T) {
		secret := &domain.Secret{
			UserID:            user.ID,
			Title:             "Old Title",
			Username:          "old",
			EncryptedPassword: "old",
		}
		require.NoError(t, secretRepo.Create(ctx, secret))

		// Update
		secret.Title = "New Title"
		err := secretRepo.Update(ctx, secret)
		require.NoError(t, err)
		assert.Equal(t, 2, secret.Version) // Should increment version

		// Verify
		found, _ := secretRepo.GetByID(ctx, secret.ID)
		assert.Equal(t, "New Title", found.Title)
		assert.Equal(t, 2, found.Version)
	})

	t.Run("DeleteSecret", func(t *testing.T) {
		secret := &domain.Secret{
			UserID:            user.ID,
			Title:             "Delete Me",
			Username:          "del",
			EncryptedPassword: "del",
		}
		require.NoError(t, secretRepo.Create(ctx, secret))

		err := secretRepo.Delete(ctx, secret.ID)
		require.NoError(t, err)

		found, err := secretRepo.GetByID(ctx, secret.ID)
		require.NoError(t, err)
		assert.Nil(t, found)
	})
}
