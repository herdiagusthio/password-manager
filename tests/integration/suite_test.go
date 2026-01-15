package integration

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testDB  *pgxpool.Pool
	connStr string
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// 1. Start Postgres Container
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpassword"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// 2. Get Connection String
	connStr, err = pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	// 3. Connect with pgxpool
	testDB, err = pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("failed to connect to db: %s", err)
	}

    // 4. Run Migrations (Manual)
	wd, _ := os.Getwd()
	projectRoot := filepath.Dir(filepath.Dir(wd))
	migrationFile := filepath.Join(projectRoot, "migrations", "000001_init.up.sql")
    
    content, err := os.ReadFile(migrationFile)
    if err != nil {
        log.Fatalf("failed to read migration file: %s", err)
    }
    
    _, err = testDB.Exec(ctx, string(content))
    if err != nil {
        log.Fatalf("failed to execute migration: %s", err)
    }

	code := m.Run()

	testDB.Close()
	_ = pgContainer.Terminate(ctx)

	os.Exit(code)
}
