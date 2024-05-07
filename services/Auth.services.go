package services

import (
	"0xKowalski1/server-hosting-web/config"
	"0xKowalski1/server-hosting-web/utils"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/google"

	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	store := NewSessionStore()

	gothic.Store = store

	goth.UseProviders(
		google.New(config.Envs.GoogleClientID, config.Envs.GoogleClientSecret, buildCallbackURL("google")),
		discord.New(config.Envs.DiscordClientID, config.Envs.DiscordClientSecret, buildCallbackURL("discord"), discord.ScopeIdentify, discord.ScopeEmail),
	)
	return &AuthService{DB: db}
}

func (service *AuthService) StoreUserSession(c echo.Context, authUser goth.User) error {
	session, err := gothic.Store.Get(c.Request(), "user_session")

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	session.Values["user"] = authUser

	err = session.Save(c.Request(), c.Response().Writer)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (service *AuthService) LogoutUser(c echo.Context) error {
	session, err := gothic.Store.Get(c.Request(), "user_session")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get session: "+err.Error())
	}

	// Delete the session values
	session.Values["user"] = nil
	session.Options.MaxAge = -1 // This deletes the cookie

	err = session.Save(c.Request(), c.Response().Writer)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save session: "+err.Error())
	}

	//
	return c.Redirect(http.StatusSeeOther, "/") // Redirect to home
}

func (service *AuthService) GetUserFromSession(c echo.Context) (goth.User, error) {
	session, err := gothic.Store.Get(c.Request(), "user_session")

	if err != nil {
		return goth.User{}, err
	}

	authUser, ok := session.Values["user"].(goth.User)
	if !ok {
		return goth.User{}, fmt.Errorf("user is not authenticated")
	}

	return authUser, nil
}

func (service *AuthService) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// User set in previous middleware
		user := utils.GetUserFromEchoContext(c)

		if user == nil {
			return c.Redirect(http.StatusTemporaryRedirect, "/login")

		}

		return next(c)
	}
}

func buildCallbackURL(provider string) string {
	return fmt.Sprintf("http://localhost:3000/auth/%s/callback", provider)
}
