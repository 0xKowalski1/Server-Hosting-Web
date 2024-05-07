package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       string `gorm:"primaryKey;"` // Set by oauth
	Email    string `gorm:"uniqueIndex"`
	Provider string

	CurrencyID uuid.UUID
	Currency   Currency

	Subscription Subscription `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
