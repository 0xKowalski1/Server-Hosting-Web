package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"0xKowalski1/server-hosting-web/db"
	"0xKowalski1/server-hosting-web/handlers"
)

func main() {
	// Connect to database
	database := db.InitDB()

	e := echo.New()

	// Static assets
	e.Static("/", "assets")

	// Middleware
	e.Use(middleware.Logger())

	// Handlers
	HomeHandler := handlers.NewHomeHandler()
	GameHandler := handlers.NewGameHandler(database)
	AuthHandler := handlers.NewAuthHandler(database)

	// Routes

	/// Home
	e.GET("/", HomeHandler.GetHome)

	/// Games
	e.GET("/games", GameHandler.GetGames)

	/// Auth
	e.GET("/login", AuthHandler.GetLogin)
	e.GET("/signup", AuthHandler.GetSignup)

	fmt.Println("Listening on :3000")
	e.Logger.Fatal(e.Start(":3000"))
}
