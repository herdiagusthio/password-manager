package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/herdiagusthio/password-manager/internal/domain"
	"github.com/herdiagusthio/password-manager/internal/usecase"
)

type AuthHandler struct {
	authUC domain.AuthUsecase
	store  *session.Store
}

func NewAuthHandler(app *fiber.App, authUC domain.AuthUsecase, store *session.Store) {
	handler := &AuthHandler{
		authUC: authUC,
		store:  store,
	}

	auth := app.Group("/auth")
	auth.Get("/login", handler.Login)
	auth.Get("/callback", handler.Callback)
	auth.Get("/logout", handler.Logout)
	auth.Get("/me", handler.Me)
}

// Login initiates the Google OIDC login flow
// @Summary Login with Google
// @Description Redirects user to Google for authentication
// @Tags Auth
// @Success 307 {string} string "Redirect to Google"
// @Router /auth/login [get]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	state := usecase.GenerateRandomState()
	
	// Store state in session to verify later (CSRF protection)
	sess, err := h.store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	sess.Set("oauthStatus", state)
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	url := h.authUC.GetLoginURL(state)
	return c.Redirect(url)
}

// Callback handles the Google OIDC callback
// @Summary Google Callback
// @Description Exchanges code for token and creates user session
// @Tags Auth
// @Param code query string true "Auth Code"
// @Param state query string true "State"
// @Success 200 {object} domain.User
// @Router /auth/callback [get]
func (h *AuthHandler) Callback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	sess, err := h.store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	} // Retrieve session

	savedState := sess.Get("oauthStatus")
	if savedState != state {
		return c.Status(fiber.StatusForbidden).SendString("Invalid state parameter")
	}

	// Remove state from session
	sess.Delete("oauthStatus")

	user, err := h.authUC.HandleCallback(c.Context(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Save user ID in session
	sess.Set("user_id", user.ID)
	sess.Set("email", user.Email)
	sess.Save()

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"user":    user,
	})
}

// Logout destroys the session
// @Summary Logout
// @Description Destroys user session
// @Tags Auth
// @Success 200 {string} string "Logged out"
// @Router /auth/logout [get]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString("Logged out")
	}
	// Destroy session
	if err := sess.Destroy(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString("Logged out")
}

// Me returns current user info
// @Summary Get Current User
// @Description Returns the authenticated user's ID and Email
// @Tags Auth
// @Success 200 {object} map[string]string
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}
	userID := sess.Get("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Not logged in"})
	}
	return c.JSON(fiber.Map{
		"user_id": userID,
		"email":   sess.Get("email"),
	})
}
