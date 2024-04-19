package handlers

import (
	"0xKowalski1/server-hosting-web/models"
	"0xKowalski1/server-hosting-web/templates"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	// Always return layout if not HTMX request, Layout will fetch the page routes SPA style
	if ctx.Request().Header.Get("HX-Request") == "" {
		var user *models.User
		userInterface := ctx.Get("user")
		if userInterface != nil {
			userConversion, ok := userInterface.(*models.User)
			if ok {
				user = userConversion
			}
		}

		return templates.Layout(t, user).Render(ctx.Request().Context(), ctx.Response().Writer)
	}

	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}
