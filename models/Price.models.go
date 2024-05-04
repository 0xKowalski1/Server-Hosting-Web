package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Price struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:uuid;primary_key;"`
	Type         string    // "Memory", "Storage", "Archive"
	PricePerUnit float32

	CurrencyID uuid.UUID
	Currency   Currency
}

// Set unique ID
func (price *Price) BeforeCreate(tx *gorm.DB) (err error) {
	if price.ID == uuid.Nil {
		price.ID = uuid.New()
	}
	return
}
