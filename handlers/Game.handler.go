package handlers

import (
	"0xKowalski1/server-hosting-web/services"
	"0xKowalski1/server-hosting-web/templates"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type GameHandler struct {
	Service *services.GameService
}

func NewGameHandler(db *gorm.DB) *GameHandler {
	return &GameHandler{
		Service: services.NewGameService(db),
	}
}

func (gh *GameHandler) GetGames(c echo.Context) error {
	games, err := gh.Service.GetGames()

	if err != nil {
		// DO something
	}

	return Render(c, 200, templates.GamesPage(games))
}
