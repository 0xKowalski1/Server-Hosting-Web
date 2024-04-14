package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/a-h/templ"

	"0xKowalski1/server-hosting-web/templates"
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

func HomeHandler(c echo.Context) error {
	return Render(c, http.StatusOK, templates.HomePage())
}

func GamesHandler(c echo.Context) error {
	// Get Games
	return Render(c, http.StatusOK, templates.GamesPage())
}

func LoginHandler(c echo.Context) error {
	// Get Games
	return Render(c, http.StatusOK, templates.LoginPage())
}

func SignupHandler(c echo.Context) error {
	// Get Games
	return Render(c, http.StatusOK, templates.SignupPage())
}

func main() {
	e := echo.New()

	e.Static("/", "assets")

	e.Use(middleware.Logger())

	e.GET("/", HomeHandler)
	e.GET("/games", GamesHandler)
	e.GET("/login", LoginHandler)
	e.GET("/signup", SignupHandler)

	fmt.Println("Listening on :3000")
	e.Logger.Fatal(e.Start(":3000"))
}
