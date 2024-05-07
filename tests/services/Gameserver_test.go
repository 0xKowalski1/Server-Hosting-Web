package services_test

import (
	"0xKowalski1/server-hosting-web/models"
	"0xKowalski1/server-hosting-web/services"
	testutils "0xKowalski1/server-hosting-web/tests/utils"
	"testing"

	"github.com/google/uuid"
)

func TestGameserverService_CreateValidGameserver(t *testing.T) {
	db := testutils.SetupDatabase()
	stripeService := services.NewStripeService(db)
	service := services.NewGameserverService(db, nil)
	currency := testutils.SetupCurrency(db)
	game := testutils.SetupGame(db)
	memoryPrice, storagePrice, archivePrice := testutils.SetupPrices(db, currency)
	user := testutils.SetupUser(db, currency)
	testutils.SetupSubscription(db, user, memoryPrice, storagePrice, archivePrice)

	newGameserver := models.Gameserver{
		Name:         "Test Server",
		MemoryLimit:  1,
		StorageLimit: 2,
		GameID:       game.ID,
		UserID:       user.ID,
	}

	gameserver, err := service.CreateGameserver(newGameserver, stripeService)

	if err != nil {
		t.Errorf("Failed to create gameserver: %v", err)
	}

	if gameserver.ID == uuid.Nil {
		t.Errorf("Failed to set gameserver ID upon creation")
	}

	// Verify gameserver is in database
	var dbGameserver models.Gameserver
	if err := db.First(&dbGameserver, gameserver.ID).Error; err != nil {
		t.Errorf("Failed to retrieve gameserver from database: %v", err)
	}
}

func TestGameserverService_CreateInvalidGameserver_NoUser(t *testing.T) {
	db := testutils.SetupDatabase()
	stripeService := services.NewStripeService(db)
	service := services.NewGameserverService(db, nil)
	game := testutils.SetupGame(db)
	currency := testutils.SetupCurrency(db)
	user := testutils.SetupUser(db, currency)

	newGameserver := models.Gameserver{
		Name:         "Test Server",
		MemoryLimit:  1,
		StorageLimit: 1,
		GameID:       game.ID,
		UserID:       user.ID,
	}

	_, err := service.CreateGameserver(newGameserver, stripeService)
	if err == nil {
		t.Error("Gameserver created!")
	}
}

func TestGameserverService_CreateInvalidGameserver_NoSubscription(t *testing.T) {
	db := testutils.SetupDatabase()
	stripeService := services.NewStripeService(db)
	service := services.NewGameserverService(db, nil)
	currency := testutils.SetupCurrency(db)
	game := testutils.SetupGame(db)
	user := testutils.SetupUser(db, currency)

	invalidGameserver := models.Gameserver{
		Name:         "Test Server",
		MemoryLimit:  1,
		StorageLimit: 1,
		GameID:       game.ID,
		UserID:       user.ID,
	}

	_, err := service.CreateGameserver(invalidGameserver, stripeService)
	if err == nil {
		t.Error("Succeeded in creating gameserver when user was not subscribed.")
	}
}

func TestGameserverService_CreateInvalidGameserver_NotEnoughResources(t *testing.T) {
	db := testutils.SetupDatabase()
	stripeService := services.NewStripeService(db)
	service := services.NewGameserverService(db, nil)
	currency := testutils.SetupCurrency(db)
	game := testutils.SetupGame(db)
	memoryPrice, storagePrice, archivePrice := testutils.SetupPrices(db, currency)
	user := testutils.SetupUser(db, currency)
	subscription := testutils.SetupSubscription(db, user, memoryPrice, storagePrice, archivePrice)

	validGameserver := models.Gameserver{
		Name:         "Test Server1",
		MemoryLimit:  subscription.MemoryGB,
		StorageLimit: subscription.StorageGB,
		GameID:       game.ID,
		UserID:       user.ID,
	}

	_, err := service.CreateGameserver(validGameserver, stripeService)
	if err != nil {
		t.Error("Failed to create valid preliminary gameserver.")
	}

	invalidGameserver := models.Gameserver{
		Name:         "Test Server2",
		MemoryLimit:  1,
		StorageLimit: 1,
		GameID:       game.ID,
		UserID:       user.ID,
	}

	_, err = service.CreateGameserver(invalidGameserver, stripeService)
	if err == nil {
		t.Error("Succeeded in creating gameserver when no resources where available.")
	}
}
