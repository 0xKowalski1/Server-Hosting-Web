package services

import (
	"0xKowalski1/server-hosting-web/config"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
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
		// Get user from session
		authUser, err := service.GetUserFromSession(c)
		if err != nil {
			log.Println("User is not authenticated:", err)
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}

		log.Printf("User is authenticated! user: %v", authUser)
		// Store user in context if needed or pass as is
		c.Set("user", authUser)
		return next(c)
	}
}

func buildCallbackURL(provider string) string {
	return fmt.Sprintf("http://localhost:3000/auth/%s/callback", provider)
}
