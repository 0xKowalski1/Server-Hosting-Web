package main

import (
	"fmt"

	"github.com/labstack/echo/v4"

	echomiddleware "github.com/labstack/echo/v4/middleware"

	"0xKowalski1/server-hosting-web/config"
	"0xKowalski1/server-hosting-web/db"
	"0xKowalski1/server-hosting-web/handlers"
	"0xKowalski1/server-hosting-web/middleware"
	"0xKowalski1/server-hosting-web/services"

	"github.com/stripe/stripe-go/v78"

	Orchestrator "0xKowalski1/container-orchestrator/api-wrapper"
)

func main() {
	// Connect to database
	database := db.InitDB()

	// Create a container orchestrator wrapper
	orchestratorWrapper := Orchestrator.NewApiWrapper("localhost") // Get me from env

	// Init Stripe
	stripe.Key = config.Envs.StripeSecretKey

	e := echo.New()

	// Static assets
	e.Static("/", "assets")

	// Services
	AuthService := services.NewAuthService(database)
	UserService := services.NewUserService(database)
	GameserverService := services.NewGameserverService(database, orchestratorWrapper)
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
	GameserverHandler := handlers.NewGameserverHandler(GameserverService, GameService, StripeService)

	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(middleware.AttachUserToContext(AuthService, UserService))

	// Configure CORS
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
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
	e.POST("/store", StoreHandler.SubmitStoreForm, middleware.RequireAuth)
	e.GET("/store/callback", StoreHandler.StripeSuccessCallback, middleware.RequireAuth)

	/// Profile
	e.GET("/profile/gameservers", GameserverHandler.GetGameservers, middleware.RequireAuth)
	e.GET("/profile/gameservers/new", GameserverHandler.NewGameserverForm, middleware.RequireAuth)
	e.POST("/profile/gameservers", GameserverHandler.CreateGameserver, middleware.RequireAuth)
	e.POST("/profile/gameservers/:id/deploy", GameserverHandler.DeployGameserver, middleware.RequireAuth, middleware.EnsureGameserverOwner(GameserverService))
	e.POST("/profile/gameservers/:id/archive", GameserverHandler.ArchiveGameserver, middleware.RequireAuth, middleware.EnsureGameserverOwner(GameserverService))
	e.POST("/profile/gameservers/:id/start", GameserverHandler.StartGameserver, middleware.RequireAuth, middleware.EnsureGameserverOwner(GameserverService))
	e.POST("/profile/gameservers/:id/stop", GameserverHandler.StopGameserver, middleware.RequireAuth, middleware.EnsureGameserverOwner(GameserverService))
	e.POST("/profile/gameservers/:id/restart", GameserverHandler.RestartGameserver, middleware.RequireAuth, middleware.EnsureGameserverOwner(GameserverService))

	fmt.Printf("Listening on :%s", config.Envs.Port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Envs.Port)))
}
