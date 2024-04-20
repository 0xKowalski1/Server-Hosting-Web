package handlers

import (
	"0xKowalski1/server-hosting-web/services"
	"0xKowalski1/server-hosting-web/templates"

	"github.com/labstack/echo/v4"
)

type GameHandler struct {
	GameService *services.GameService
}

func NewGameHandler(gameService *services.GameService) *GameHandler {
	return &GameHandler{
		GameService: gameService,
	}
}

func (gh *GameHandler) GetGames(c echo.Context) error {
	searchQuery := c.QueryParam("search")

	games, err := gh.GameService.GetGames(searchQuery)

	if err != nil {
		// DO something
		return err
	}

	if c.Request().Header.Get("X-Partial-Content") == "true" {
		return Render(c, 200, templates.GamesList(games))
	} else {
		return Render(c, 200, templates.GamesPage(games))
	}
}
