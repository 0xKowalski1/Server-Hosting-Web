package handlers

import (
	"0xKowalski1/server-hosting-web/models"
	"0xKowalski1/server-hosting-web/services"
	"0xKowalski1/server-hosting-web/templates"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type GameserverHandler struct {
	Service *services.GameserverService
}

func NewGameserverHandler(db *gorm.DB) *GameserverHandler {
	return &GameserverHandler{
		Service: services.NewGameserverService(db),
	}
}

func (gh *GameserverHandler) NewGameserverForm(c echo.Context) error {
	formData := models.GameserverFormData{
		Name: "jeff",
	}
	return Render(c, 200, templates.GameserverForm(formData))
}

func (gh *GameserverHandler) GetGameservers(c echo.Context) error {
	gameservers, err := gh.Service.GetGameservers()

	if err != nil {
		// Do something
	}

	return Render(c, 200, templates.GameserversPage(gameservers))
}

func (gh *GameserverHandler) CreateGameserver(c echo.Context) error {
	newGameserver := models.Gameserver{
		Name: c.FormValue("name"),
	}
	_, err := gh.Service.CreateGameserver(newGameserver)

	if err != nil {
		// Should send back errors for new gameserver form
		return err
	}

	return nil
}
