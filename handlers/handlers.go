package handlers

import (
	"0xKowalski1/server-hosting-web/templates"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	// Always return layout if not HTMX request, Layout will fetch the page routes SPA style
	if ctx.Request().Header.Get("HX-Request") == "" {
		return templates.Layout().Render(ctx.Request().Context(), ctx.Response().Writer)
	}

	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}
