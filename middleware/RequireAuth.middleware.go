package middleware

import (
	"0xKowalski1/server-hosting-web/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := utils.GetUserFromEchoContext(c)

		if user == nil {
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}

		return next(c)
	}
}
