package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Currency struct {
	gorm.Model
	ID     uuid.UUID `gorm:"type:uuid;primary_key;"`
	Code   string    // ISO currency code, e.g., USD, EUR
	Symbol string    // $
}

// Set unique ID
func (currency *Currency) BeforeCreate(tx *gorm.DB) (err error) {
	if currency.ID == uuid.Nil {
		currency.ID = uuid.New()
	}
	return
}
