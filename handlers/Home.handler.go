package handlers

import (
	"0xKowalski1/server-hosting-web/templates"

	"github.com/labstack/echo/v4"
)

type HomeHandler struct {
}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (hh *HomeHandler) GetHome(c echo.Context) error {
	return Render(c, 200, templates.HomePage())
}
