package handlers

import (
	"0xKowalski1/server-hosting-web/templates"

	"github.com/labstack/echo/v4"
)

type SupportHandler struct {
}

func NewSupportHandler() *SupportHandler {
	return &SupportHandler{}
}

func (sh *SupportHandler) GetSupport(c echo.Context) error {
	return Render(c, 200, templates.SupportPage())
}
