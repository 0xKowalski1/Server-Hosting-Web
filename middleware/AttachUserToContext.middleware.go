package middleware

import (
	"0xKowalski1/server-hosting-web/services"

	"github.com/labstack/echo/v4"
)

func AttachUserToContext(authService *services.AuthService, userService *services.UserService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authUser, err := authService.GetUserFromSession(c)
			if err != nil {
				c.Set("user", nil)
				return next(c)
			}

			dbUser, err := userService.GetUser(authUser.UserID)
			if err != nil {
				c.Set("user", nil)
				return next(c)
			}

			c.Set("user", dbUser)
			return next(c)
		}
	}
}
