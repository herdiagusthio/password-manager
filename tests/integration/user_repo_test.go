package integration

import (
	"context"
	"testing"
	"time"

	"github.com/herdiagusthio/password-manager/internal/domain"
	"github.com/herdiagusthio/password-manager/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepo(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: database not initialized")
	}

	repo := postgres.NewUserRepository(testDB)
	ctx := context.Background()

	t.Run("CreateUser", func(t *testing.T) {
		email := "test@example.com"
		user := &domain.User{
			Email: email,
		}

		err := repo.Create(ctx, user)
		require.NoError(t, err)
		assert.NotEmpty(t, user.ID)
		assert.Equal(t, email, user.Email)
		assert.WithinDuration(t, time.Now(), user.CreatedAt, 2*time.Second)
	})

	t.Run("GetUserByEmail", func(t *testing.T) {
		email := "findme@example.com"
		user := &domain.User{Email: email}
		require.NoError(t, repo.Create(ctx, user))

		found, err := repo.GetByEmail(ctx, email)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, user.ID, found.ID)
		assert.Equal(t, email, found.Email)
	})

	t.Run("GetNonExistentUser", func(t *testing.T) {
		found, err := repo.GetByEmail(ctx, "ghost@example.com")
		require.NoError(t, err)
		assert.Nil(t, found)
	})
}
