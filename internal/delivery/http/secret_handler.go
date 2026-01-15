package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/herdiagusthio/password-manager/internal/domain"
)

type SecretHandler struct {
	usecase domain.SecretUsecase
	store   *session.Store
}

func NewSecretHandler(app *fiber.App, uc domain.SecretUsecase, store *session.Store) {
	h := &SecretHandler{
		usecase: uc,
		store:   store,
	}

	api := app.Group("/api", h.requireAuth)
	api.Post("/secrets", h.Create)
	api.Get("/secrets", h.List)
	api.Get("/secrets/:id", h.Get)
	api.Put("/secrets/:id", h.Update)
	api.Delete("/secrets/:id", h.Delete)
}

func (h *SecretHandler) requireAuth(c *fiber.Ctx) error {
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

// Create creates a new secret
// @Summary Create Secret
// @Description Create a new encrypted secret
// @Tags Secrets
// @Accept json
// @Produce json
// @Param secret body object true "Secret Data"
// @Success 201 {object} domain.Secret
// @Router /api/secrets [post]
func (h *SecretHandler) Create(c *fiber.Ctx) error {
	type Request struct {
		Title    string                 `json:"title"`
		Username string                 `json:"username"`
		Password string                 `json:"password"`
		Metadata map[string]interface{} `json:"metadata"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userID := c.Locals("user_id").(string)
	secret := &domain.Secret{
		UserID:   userID,
		Title:    req.Title,
		Username: req.Username,
		Password: req.Password,
		Metadata: req.Metadata,
	}

	if err := h.usecase.CreateSecret(c.Context(), secret); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(secret)
}

// List returns all secrets for the user
// @Summary List Secrets
// @Description Get all secrets (without passwords)
// @Tags Secrets
// @Produce json
// @Success 200 {array} domain.Secret
// @Router /api/secrets [get]
func (h *SecretHandler) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	secrets, err := h.usecase.ListSecrets(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(secrets)
}

// Get returns a single secret (decrypted)
// @Summary Get Secret
// @Description Get a secret by ID with decrypted password
// @Tags Secrets
// @Produce json
// @Param id path string true "Secret ID"
// @Success 200 {object} domain.Secret
// @Router /api/secrets/{id} [get]
func (h *SecretHandler) Get(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	id := c.Params("id")

	secret, err := h.usecase.GetSecret(c.Context(), id, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if secret == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.JSON(secret)
}

// Update modifies an existing secret
// @Summary Update Secret
// @Description Update secret details
// @Tags Secrets
// @Accept json
// @Produce json
// @Param id path string true "Secret ID"
// @Param secret body object true "Secret Data"
// @Success 200 {object} domain.Secret
// @Router /api/secrets/{id} [put]
func (h *SecretHandler) Update(c *fiber.Ctx) error {
	// Simplied update...
	type Request struct {
		Title    string                 `json:"title"`
		Username string                 `json:"username"`
		Password string                 `json:"password"`
		Metadata map[string]interface{} `json:"metadata"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userID := c.Locals("user_id").(string)
	id := c.Params("id")

	secret := &domain.Secret{
		ID:       id,
		UserID:   userID,
		Title:    req.Title,
		Username: req.Username,
		Password: req.Password,
		Metadata: req.Metadata,
	}

	if err := h.usecase.UpdateSecret(c.Context(), secret); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(secret)
}

// Delete removes a secret
// @Summary Delete Secret
// @Description Remove a secret by ID
// @Tags Secrets
// @Param id path string true "Secret ID"
// @Success 204 "No Content"
// @Router /api/secrets/{id} [delete]
func (h *SecretHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	id := c.Params("id")

	if err := h.usecase.DeleteSecret(c.Context(), id, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
