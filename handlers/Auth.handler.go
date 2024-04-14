package handlers

import (
	//	"0xKowalski1/server-hosting-web/services"
	"0xKowalski1/server-hosting-web/templates"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AuthHandler struct {
	// service *services.authservice
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		//	service: services.newgameservice(db),
	}
}

func (ah *AuthHandler) GetLogin(c echo.Context) error {
	return Render(c, 200, templates.LoginPage())
}

func (ah *AuthHandler) GetSignup(c echo.Context) error {
	return Render(c, 200, templates.SignupPage())
}
