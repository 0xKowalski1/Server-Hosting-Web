package services

import (
	"0xKowalski1/server-hosting-web/models"

	"gorm.io/gorm"
)

type GameService struct {
	DB *gorm.DB
}

func NewGameService(db *gorm.DB) *GameService {
	return &GameService{DB: db}
}

func (service *GameService) GetGames() ([]models.Game, error) {
	var games []models.Game
	result := service.DB.Find(&games)
	return games, result.Error
}
