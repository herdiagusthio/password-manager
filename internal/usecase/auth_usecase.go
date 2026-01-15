package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/herdiagusthio/password-manager/config"
	"github.com/herdiagusthio/password-manager/internal/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type authUsecase struct {
	oauthConfig *oauth2.Config
	userRepo    domain.AuthRepository
	cfg         *config.Config
}

func NewAuthUsecase(cfg *config.Config, userRepo domain.AuthRepository) domain.AuthUsecase {
	conf := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &authUsecase{
		oauthConfig: conf,
		userRepo:    userRepo,
		cfg:         cfg,
	}
}

func (u *authUsecase) GetLoginURL(state string) string {
	return u.oauthConfig.AuthCodeURL(state)
}

func (u *authUsecase) HandleCallback(ctx context.Context, code string) (*domain.User, error) {
	// 1. Exchange code for token
	token, err := u.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// 2. Fetch user info from Google
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google returned non-200 status: %d", resp.StatusCode)
	}

	userInfoBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 3. Parse user info
	var googleUser struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(userInfoBytes, &googleUser); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	// 4. Find or Create User in DB
	user, err := u.userRepo.GetByEmail(ctx, googleUser.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		// Create new user
		user = &domain.User{
			Email: googleUser.Email,
		}
		if err := u.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}

// GenerateRandomState generates a random state string for CSRF protection
func GenerateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
