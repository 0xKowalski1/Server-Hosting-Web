package handlers

import (
	"0xKowalski1/server-hosting-web/models"
	"0xKowalski1/server-hosting-web/services"
	"0xKowalski1/server-hosting-web/templates"
	"0xKowalski1/server-hosting-web/utils"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
)

type GameserverHandler struct {
	GameserverService *services.GameserverService
	GameService       *services.GameService
	StripeService     *services.StripeService
}

func NewGameserverHandler(gameserverService *services.GameserverService, gameService *services.GameService, stripeService *services.StripeService) *GameserverHandler {
	return &GameserverHandler{
		GameserverService: gameserverService,
		GameService:       gameService,
		StripeService:     stripeService,
	}
}

func (gh *GameserverHandler) NewGameserverForm(c echo.Context) error {
	games, err := gh.GameService.GetGames("")

	if err != nil {
		// DO something
		return err
	}

	formData := models.GameserverFormData{
		Name:   "mc",
		GameID: games[0].ID.String(),
	}
	return Render(c, 200, templates.GameserverForm(formData, games))
}

func (gh *GameserverHandler) GetGameservers(c echo.Context) error {
	user := utils.GetUserFromEchoContext(c)

	gameservers, err := gh.GameserverService.GetGameservers(user.ID)

	if err != nil {
		// Do something
	}

	usedMemory, usedStorage := gh.GameserverService.GetUsedResources(gameservers)

	var usedMemoryPercentage, usedStoragePercentage float32
	// Handle case where subscription is 0, which would cause a divide by 0
	if user.Subscription.MemoryGB > 0 {
		usedMemoryPercentage = (float32(usedMemory) / float32(user.Subscription.MemoryGB)) * 100
	}
	if user.Subscription.StorageGB > 0 {
		usedStoragePercentage = (float32(usedStorage) / float32(user.Subscription.StorageGB)) * 100
	}

	return Render(c, 200, templates.GameserversPage(gameservers, usedMemory, usedStorage, usedMemoryPercentage, usedStoragePercentage, user))
}

func (gh *GameserverHandler) CreateGameserver(c echo.Context) error {
	user := utils.GetUserFromEchoContext(c)

	game, err := gh.GameService.GetGameByID(c.FormValue("game"))
	if err != nil {
		//Do something
		log.Printf("Error getting game by id: %v", err)
		return err
	}

	memoryLimit, err := strconv.Atoi(c.FormValue("memory"))
	if err != nil {
		log.Printf("Error converting memory to int: %v", err)
		return err
	}
	storageLimit, err := strconv.Atoi(c.FormValue("storage"))
	if err != nil {
		log.Printf("Error converting storage to int: %v", err)
		return err
	}

	newGameserver := models.Gameserver{
		Name:         c.FormValue("name"),
		GameID:       game.ID,
		MemoryLimit:  memoryLimit,
		StorageLimit: storageLimit,
		UserID:       user.ID,
		Status:       models.GameserverStatusUndeployed,
	}

	_, err = gh.GameserverService.CreateGameserver(newGameserver, gh.StripeService)

	if err != nil {
		// Should send back errors for new gameserver form
		log.Printf("Error creating new gameserver: %v", err)
		return err
	}

	c.Response().Header().Set("HX-Replace-Url", "/profile/gameservers")
	return gh.GetGameservers(c)
}

func (gh *GameserverHandler) DeployGameserver(c echo.Context) error {
	gameserver := utils.GetGameserverFromEchoContext(c)
	if gameserver == nil {
		// 404
	}

	gameserver, err := gh.GameserverService.DeployGameserver(gameserver)
	if err != nil {
		log.Printf("Error deploying gameserver at ID - %s: %v", gameserver.ID.String(), err)
		return err
	}

	return Render(c, 200, templates.GameserverCard(*gameserver))
}

func (gh *GameserverHandler) ArchiveGameserver(c echo.Context) error {
	gameserver := utils.GetGameserverFromEchoContext(c)
	if gameserver == nil {
		// 404
	}

	gameserver, err := gh.GameserverService.ArchiveGameserver(gameserver)
	if err != nil {
		log.Printf("Error archiving gameserver at ID - %s: %v", gameserver.ID.String(), err)
		return err
	}

	return Render(c, 200, templates.GameserverCard(*gameserver))
}

func (gh *GameserverHandler) StartGameserver(c echo.Context) error {
	gameserver := utils.GetGameserverFromEchoContext(c)
	if gameserver == nil {
		// 404
	}

	gameserver, err := gh.GameserverService.StartGameserver(gameserver)
	if err != nil {
		log.Printf("Error starting gameserver at ID - %s: %v", gameserver.ID.String(), err)
		return err
	}

	return Render(c, 200, templates.GameserverCard(*gameserver))
}

func (gh *GameserverHandler) StopGameserver(c echo.Context) error {
	gameserver := utils.GetGameserverFromEchoContext(c)
	if gameserver == nil {
		// 404
	}

	gameserver, err := gh.GameserverService.StopGameserver(gameserver)
	if err != nil {
		log.Printf("Error stopping gameserver at ID - %s: %v", gameserver.ID.String(), err)
		return err
	}

	return Render(c, 200, templates.GameserverCard(*gameserver))
}

func (gh *GameserverHandler) RestartGameserver(c echo.Context) error {
	gameserver := utils.GetGameserverFromEchoContext(c)
	if gameserver == nil {
		// 404
	}

	gameserver, err := gh.GameserverService.StopGameserver(gameserver)
	if err != nil {
		log.Printf("Error stopping gameserver at ID - %s: %v", gameserver.ID.String(), err)
		return err
	}

	gameserver, err = gh.GameserverService.StartGameserver(gameserver)
	if err != nil {
		log.Printf("Error starting gameserver at ID - %s: %v", gameserver.ID.String(), err)
		return err
	}

	return Render(c, 200, templates.GameserverCard(*gameserver))
}
