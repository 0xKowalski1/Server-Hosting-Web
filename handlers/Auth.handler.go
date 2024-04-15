package handlers

import (
	"0xKowalski1/server-hosting-web/models"
	"0xKowalski1/server-hosting-web/services"
	"0xKowalski1/server-hosting-web/templates"
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

type AuthHandler struct {
	authService *services.AuthService
	userService *services.UserService
}

func NewAuthHandler(authService *services.AuthService, userService *services.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

func (ah *AuthHandler) GetLogin(c echo.Context) error {
	return Render(c, 200, templates.LoginPage())
}

func (ah *AuthHandler) BeginAuth(c echo.Context) error {
	w := c.Response().Writer
	r := c.Request()

	// Handle no provider
	provider := c.Param("provider")
	if provider == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Provider not specified")
	}

	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	if _, err := gothic.CompleteUserAuth(w, r); err == nil {
		// Handle user already logged in
		log.Println("User already logged in.")
		return Render(c, 200, templates.HomePage())
	} else {
		log.Println("Beginning auth process")
		// Start the authentication process
		gothic.BeginAuthHandler(w, r)
	}

	return nil
}

func (ah *AuthHandler) AuthCallback(c echo.Context) error {
	w := c.Response().Writer
	r := c.Request()

	// Handle no provider
	provider := c.Param("provider")
	if provider == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Provider not specified")
	}

	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	// Complete the user auth and handle the callback from the OAuth provider
	authUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Println("Authentication failed:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Authentication failed")
	}

	_, err = ah.userService.FindOrCreateUser(models.User{Provider: provider, ID: authUser.UserID, Email: authUser.Email})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create or find user")
	}

	// Store the authenticated user session
	if err := ah.authService.StoreUserSession(c, authUser); err != nil {
		log.Printf("Error storing user session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to store user session")
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
