package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       string // Set by oauth
	Email    string `gorm:"uniqueIndex"`
	Provider string

	CurrencyID uuid.UUID
	Currency   Currency
}
