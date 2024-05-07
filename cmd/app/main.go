package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"

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

	// Services
	AuthService := services.NewAuthService(database)
	UserService := services.NewUserService(database)
	GameserverService := services.NewGameserverService(database)
	GameService := services.NewGameService(database)
	CurrencyService := services.NewCurrencyService(database)
	PriceService := services.NewPriceService(database)
	StripeService := services.NewStripeService(database)

	// Handlers
	HomeHandler := handlers.NewHomeHandler()
	GameHandler := handlers.NewGameHandler(GameService)
	SupportHandler := handlers.NewSupportHandler()
	StoreHandler := handlers.NewStoreHandler(CurrencyService, PriceService, StripeService)
	AuthHandler := handlers.NewAuthHandler(AuthService, UserService)
	GameserverHandler := handlers.NewGameserverHandler(GameserverService, GameService)

	// Middleware
	//e.Use(middleware.Logger())
	e.Use(AttachUserToContext(AuthService, UserService))

	// Configure CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Routes
	/// Home
	e.GET("/", HomeHandler.GetHome)

	/// Games
	e.GET("/games", GameHandler.GetGames)

	/// Auth
	e.GET("/login", AuthHandler.GetLogin)
	e.GET("/logout", AuthHandler.PostLogout) // Does not require actual auth, all it does is deletes session cookies
	e.GET("/auth/:provider", AuthHandler.BeginAuth)
	e.GET("/auth/:provider/callback", AuthHandler.AuthCallback)

	// Support
	e.GET("/support", SupportHandler.GetSupport)

	// Store
	e.GET("/store", StoreHandler.GetStore)
	e.GET("/store/guided", StoreHandler.GetGuidedStoreFlow)
	e.GET("/store/advanced", StoreHandler.GetAdvancedStoreFlow)
	e.POST("/store", StoreHandler.SubmitStoreForm, AuthService.RequireAuth)
	e.GET("/store/callback", StoreHandler.StripeSuccessCallback, AuthService.RequireAuth)

	/// Profile
	e.GET("/profile/gameservers", GameserverHandler.GetGameservers, AuthService.RequireAuth)
	e.GET("/profile/gameservers/new", GameserverHandler.NewGameserverForm, AuthService.RequireAuth)
	e.POST("/profile/gameservers", GameserverHandler.CreateGameserver, AuthService.RequireAuth)
	e.POST("/profile/gameservers/:id/deploy", GameserverHandler.DeployGameserver, AuthService.RequireAuth)

	fmt.Printf("Listening on :%s", config.Envs.Port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Envs.Port)))
}

// Should be in a middleware dir
func AttachUserToContext(authService *services.AuthService, userService *services.UserService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authUser, err := authService.GetUserFromSession(c)
			if err != nil {
				log.Printf("Failed to get user from session: %v", err)
				c.Set("user", nil)
				return next(c)
			}

			dbUser, err := userService.GetUser(authUser.UserID)
			if err != nil {
				log.Printf("Failed to get user from DB: %v", err)
				c.Set("user", nil)
				return next(c)
			}

			c.Set("user", dbUser)
			return next(c)
		}
	}
}
