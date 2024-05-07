package utils

import (
	"0xKowalski1/server-hosting-web/models"

	"github.com/labstack/echo/v4"
)

func GetUserFromEchoContext(c echo.Context) *models.User {
	var user *models.User
	userInterface := c.Get("user")
	if userInterface != nil {
		userConversion, ok := userInterface.(*models.User)
		if ok {
			user = userConversion
		}
	}

	return user
}
