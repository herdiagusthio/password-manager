package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/herdiagusthio/password-manager/internal/domain"
)

type UIHandler struct {
	secretUC domain.SecretUsecase
	store    *session.Store
}

func NewUIHandler(app *fiber.App, secretUC domain.SecretUsecase, store *session.Store) {
	h := &UIHandler{
		secretUC: secretUC,
		store:    store,
	}

	app.Get("/", h.Landing)
	app.Get("/login", h.LoginPage)
	app.Get("/dashboard", h.requireAuth, h.Dashboard)
}

func (h *UIHandler) requireAuth(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		return c.Redirect("/login")
	}
	userID := sess.Get("user_id")
	if userID == nil {
		return c.Redirect("/login")
	}
	c.Locals("user_id", userID)
	c.Locals("email", sess.Get("email"))
	return c.Next()
}

func (h *UIHandler) Landing(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err == nil && sess.Get("user_id") != nil {
		return c.Redirect("/dashboard")
	}
	return c.Redirect("/login")
}

func (h *UIHandler) LoginPage(c *fiber.Ctx) error {
	return c.Render("auth/login", fiber.Map{
		"Authenticated": false,
	}, "layouts/main")
}

func (h *UIHandler) Dashboard(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	email := c.Locals("email").(string)

	secrets, err := h.secretUC.ListSecrets(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error fetching secrets")
	}

	return c.Render("dashboard/index", fiber.Map{
		"Authenticated": true,
		"UserEmail":     email,
		"Secrets":       secrets,
	}, "layouts/main")
}
