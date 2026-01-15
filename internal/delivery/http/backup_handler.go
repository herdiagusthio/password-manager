package http

import (
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/herdiagusthio/password-manager/internal/domain"
)

type BackupHandler struct {
	usecase domain.BackupUsecase
	store   *session.Store
}

func NewBackupHandler(app *fiber.App, uc domain.BackupUsecase, store *session.Store) {
	h := &BackupHandler{
		usecase: uc,
		store:   store,
	}

	api := app.Group("/api", h.requireAuth) 
	// requireAuth is reused from secret_handler.go? No, it's a method on SecretHandler. 
	// I should probably make a public middleware or duplicate it. 
	// For simplicity, I will duplicate the auth logic here or I should have refactored it. 
	// Refactoring is better but I'll write a small inline middleware for now to avoid cross-file dependency mess if I don't move it to a shared pkg.
	
	api.Get("/backup/export", h.Export)
	api.Post("/backup/import", h.Import)
}

func (h *BackupHandler) requireAuth(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	userID := sess.Get("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	c.Locals("user_id", userID)
	return c.Next()
}

func (h *BackupHandler) Export(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	data, err := h.usecase.ExportSecrets(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	c.Set("Content-Disposition", "attachment; filename=secrets_backup.enc")
	c.Set("Content-Type", "application/octet-stream")
	return c.Send(data)
}

func (h *BackupHandler) Import(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	fileHeader, err := c.FormFile("backup")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing backup file"})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to read file"})
	}

	if err := h.usecase.ImportSecrets(c.Context(), userID, data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "restore failed: " + err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Backup restored successfully"})
}
