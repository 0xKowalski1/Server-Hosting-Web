package services

import (
	"0xKowalski1/server-hosting-web/models"
	"strings"

	"gorm.io/gorm"
)

type GameService struct {
	DB *gorm.DB
}

func NewGameService(db *gorm.DB) *GameService {
	return &GameService{DB: db}
}

// GetGames retrieves games from the database, optionally filtering them based on a search term.
func (service *GameService) GetGames(searchQuery string) ([]models.Game, error) {
	var games []models.Game
	query := service.DB.Model(&models.Game{})

	// If a search query is provided, use it to filter the results
	if searchQuery != "" {
		searchQuery = "%" + strings.ToLower(searchQuery) + "%" // Prepare the search query for case-insensitive matching
		query = query.Where("lower(name) LIKE ?", searchQuery)
	}

	result := query.Find(&games)
	// handle errors
	if result.Error != nil {
		return nil, result.Error
	}

	return games, nil
}
