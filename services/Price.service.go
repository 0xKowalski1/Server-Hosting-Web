package services

import (
	"0xKowalski1/server-hosting-web/models"
	"gorm.io/gorm"
)

type PriceService struct {
	DB *gorm.DB
}

func NewPriceService(db *gorm.DB) *PriceService {
	return &PriceService{DB: db}
}

// Gets all prices for a specific currency
func (service *PriceService) GetPrices(currency *models.Currency) ([]models.Price, error) {
	var prices []models.Price

	result := service.DB.Where("currency_id = ?", currency.ID).Find(&prices)
	if result.Error != nil {
		return nil, result.Error
	}

	return prices, nil
}
