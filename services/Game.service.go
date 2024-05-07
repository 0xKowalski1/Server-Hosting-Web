package services

import (
	"0xKowalski1/server-hosting-web/models"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GameService struct {
	DB *gorm.DB
}

func NewGameService(db *gorm.DB) *GameService {
	return &GameService{DB: db}
}

func (service *GameService) GetGames(searchQuery string) ([]models.Game, error) {
	var games []models.Game
	query := service.DB.Model(&models.Game{})

	if searchQuery != "" {
		searchQuery = "%" + strings.ToLower(searchQuery) + "%" // Prepare the search query for case-insensitive matching
		query = query.Where("lower(name) LIKE ?", searchQuery)
	}

	result := query.Find(&games)
	if result.Error != nil {
		return nil, fmt.Errorf("Failed to query games in database: %v", result.Error)
	}

	return games, nil
}

// Should take gameID as UUID
func (service *GameService) GetGameByID(gameID string) (*models.Game, error) {
	var game *models.Game

	id, err := uuid.Parse(gameID)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse uuid from gameID: %s due to error: %v", gameID, err)
	}

	result := service.DB.First(&game, "id = ?", id)
	if result.Error != nil {
		return nil, fmt.Errorf("Failed to find game at ID: %s due to error: %v", gameID, result.Error)
	}

	return game, nil
}
