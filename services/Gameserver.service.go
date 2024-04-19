package services

import (
	"0xKowalski1/server-hosting-web/models"
	"log"

	"gorm.io/gorm"
)

type GameserverService struct {
	DB *gorm.DB
}

func NewGameserverService(db *gorm.DB) *GameserverService {
	return &GameserverService{DB: db}
}

func (service *GameserverService) CreateGameserver(newGameserver models.Gameserver) (*models.Gameserver, error) {
	result := service.DB.Create(&newGameserver)
	if result.Error != nil {
		log.Fatalf("Failed to create gameserver: %v", result.Error)
		return nil, result.Error
	}

	return &newGameserver, nil
}

func (service *GameserverService) GetGameservers() ([]models.Gameserver, error) {
	var gameservers []models.Gameserver

	result := service.DB.Find(&gameservers)
	if result.Error != nil {
		return nil, result.Error
	}

	return gameservers, nil
}
