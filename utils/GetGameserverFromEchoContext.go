package utils

import (
	"0xKowalski1/server-hosting-web/models"

	"github.com/labstack/echo/v4"
)

func GetGameserverFromEchoContext(c echo.Context) *models.Gameserver {
	var gameserver *models.Gameserver
	gameserverInterface := c.Get("gameserver")
	if gameserverInterface != nil {
		gameserverConversion, ok := gameserverInterface.(*models.Gameserver)
		if ok {
			gameserver = gameserverConversion
		}
	}

	return gameserver
}
