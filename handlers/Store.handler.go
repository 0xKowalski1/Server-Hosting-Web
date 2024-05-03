package handlers

import (
	"0xKowalski1/server-hosting-web/templates"

	"github.com/labstack/echo/v4"
)

type StoreHandler struct {
}

func NewStoreHandler() *StoreHandler {
	return &StoreHandler{}
}

func (sh *StoreHandler) GetStore(c echo.Context) error {
	return Render(c, 200, templates.StorePage())
}

func (sh *StoreHandler) GetGuidedStoreFlow(c echo.Context) error {
	return Render(c, 200, templates.GuidedStoreFlow())
}

func (sh *StoreHandler) GetAdvancedStoreFlow(c echo.Context) error {
	return Render(c, 200, templates.AdvancedStoreFlow())
}

func (sh *StoreHandler) GetCheckout(c echo.Context) error {
	return Render(c, 200, templates.Checkout())

}
