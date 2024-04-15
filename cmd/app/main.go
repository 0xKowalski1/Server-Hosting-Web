package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	//	"github.com/labstack/echo/v4/middleware"

	"0xKowalski1/server-hosting-web/config"
	"0xKowalski1/server-hosting-web/db"
	"0xKowalski1/server-hosting-web/handlers"
	"0xKowalski1/server-hosting-web/services"
)

func main() {
	// Connect to database
	database := db.InitDB()

	e := echo.New()

	// Static assets
	e.Static("/", "assets")

	// Middleware
	//e.Use(middleware.Logger())

	// Services
	AuthService := services.NewAuthService(database)
	UserService := services.NewUserService(database)

	// Handlers
	HomeHandler := handlers.NewHomeHandler()
	GameHandler := handlers.NewGameHandler(database)
	AuthHandler := handlers.NewAuthHandler(AuthService, UserService)

	// Routes

	/// Home
	e.GET("/", HomeHandler.GetHome)

	/// Games
	e.GET("/games", GameHandler.GetGames, AuthService.RequireAuth)

	/// Auth
	e.GET("/login", AuthHandler.GetLogin)
	e.GET("/signup", AuthHandler.GetSignup)
	e.GET("/auth/:provider", AuthHandler.BeginAuth)
	e.GET("/auth/:provider/callback", AuthHandler.AuthCallback)

	fmt.Printf("Listening on :%s", config.Envs.Port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Envs.Port)))
}
