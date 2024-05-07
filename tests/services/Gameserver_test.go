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
	service := services.NewGameserverService(db, nil)
	//currency := testutils.SetupCurrency(db)
	game := testutils.SetupGame(db)
	//user := testutils.SetupUser(db, currency)

	newGameserver := models.Gameserver{
		Name:         "Test Server",
		MemoryLimit:  1,
		StorageLimit: 2,
		GameID:       game.ID,
	}

	gameserver, err := service.CreateGameserver(newGameserver)

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

	if dbGameserver.Name != gameserver.Name {
		t.Errorf("Retrieved gameserver name '%v' does not match expected name '%v'", dbGameserver.Name, gameserver.Name)
	}
}
