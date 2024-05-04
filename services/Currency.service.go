package services

import (
	"0xKowalski1/server-hosting-web/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CurrencyService struct {
	DB *gorm.DB
}

func NewCurrencyService(db *gorm.DB) *CurrencyService {
	return &CurrencyService{DB: db}
}

func (service *CurrencyService) GetCurrencies() ([]models.Currency, error) {
	var currencies []models.Currency

	result := service.DB.Find(&currencies)
	if result.Error != nil {
		return nil, result.Error
	}

	return currencies, nil
}

func (service *CurrencyService) GetCurrencyById(currencyID uuid.UUID) (models.Currency, error) {
	var currency models.Currency

	result := service.DB.First(&currency, "id = ?", currencyID)
	if result.Error != nil {
		return currency, result.Error
	}

	return currency, nil
}

func (service *CurrencyService) GetDefaultCurrency() (models.Currency, error) {
	var currency models.Currency

	result := service.DB.First(&currency, "code = ?", "USD")
	if result.Error != nil {
		return currency, result.Error
	}

	return currency, nil
}
