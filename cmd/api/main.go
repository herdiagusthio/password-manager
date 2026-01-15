package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis/v3"
	"github.com/gofiber/template/html/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/herdiagusthio/password-manager/config"
	authHttp "github.com/herdiagusthio/password-manager/internal/delivery/http"
	postgresRepo "github.com/herdiagusthio/password-manager/internal/repository/postgres"
	"github.com/herdiagusthio/password-manager/internal/usecase"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Database Connection (Postgres)
	dbPool, err := pgxpool.New(context.Background(), cfg.DBSource)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}
	log.Println("Connected to Database")

	// 3. Redis Configuration for Session Store
	// Note: We are using "github.com/gofiber/storage/redis/v3" for Fiber Session
	// It requires a slightly different setup than standard go-redis client if we use it directly for storage
	// But `session.New` takes a Config.
	redisStorage := redis.New(redis.Config{
		URL: "redis://" + cfg.RedisAddr,
	})
	
	sessionStore := session.New(session.Config{
		Storage: redisStorage,
		Expiration: 24 * time.Hour, 
		KeyLookup: "cookie:session_id",
	})

	// 4. Initialize Fiber App
	engine := html.New("./views", ".html")
	
	app := fiber.New(fiber.Config{
		AppName: "Password Manager API",
		Views:   engine,
	})
	app.Use(logger.New())
	app.Use(recover.New())

	// 5. Dependency Injection
	// Repositories
	userRepo := postgresRepo.NewUserRepository(dbPool)
	secretRepo := postgresRepo.NewSecretRepository(dbPool)

	// Usecases
	authUC := usecase.NewAuthUsecase(&cfg, userRepo)
	secretUC := usecase.NewSecretUsecase(secretRepo, &cfg)
	backupUC := usecase.NewBackupUsecase(secretRepo, &cfg)

	// Handlers
	authHttp.NewAuthHandler(app, authUC, sessionStore)
	authHttp.NewSecretHandler(app, secretUC, sessionStore)
	authHttp.NewBackupHandler(app, backupUC, sessionStore)
	authHttp.NewUIHandler(app, secretUC, sessionStore)

	// Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// 6. Graceful Shutdown & Server Start
	go func() {
		if err := app.Listen(cfg.ServerPort); err != nil {
			log.Printf("Server Listen Error: %v", err)
		}
	}()

	log.Printf("Server started on %s", cfg.ServerPort)

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down...")
	app.Shutdown()
}
