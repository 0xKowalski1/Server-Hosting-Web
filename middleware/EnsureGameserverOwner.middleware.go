package middleware

import (
	"0xKowalski1/server-hosting-web/services"
	"0xKowalski1/server-hosting-web/utils"
	"log"

	"github.com/labstack/echo/v4"
)

func EnsureGameserverOwner(gameserverService *services.GameserverService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			gameserverID := c.Param("id")
			gameserver, err := gameserverService.GetGameserverByID(gameserverID)
			if err != nil {
				log.Printf("Error finding gameserver at ID - %s: %v", gameserverID, err)
				c.Set("gameserver", nil)
				return next(c)
			}

			user := utils.GetUserFromEchoContext(c)
			if gameserver.UserID != user.ID {
				log.Printf("User does not own gameserver at ID - %s: %v", gameserverID, err)
				c.Set("gameserver", nil)
				return next(c)
			}

			c.Set("gameserver", gameserver)
			return next(c)
		}
	}
}

/*


 */
