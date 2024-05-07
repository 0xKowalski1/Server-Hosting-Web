package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GameserverFormData struct {
	Name   string
	GameID string
}

type Gameserver struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;primary_key;" validate:"required,uuid"`
	Name string    `validate:"required,min=3,max=100"`

	// Belongs To User
	UserID string `validate:"required"`

	// Belongs to Game
	GameID uuid.UUID `validate:"required,uuid"`
	Game   Game

	// Resource Limits
	StorageLimit int `validate:"required,min=1,max=128"`
	MemoryLimit  int `validate:"required,min=1,max=32"`
}

// Set unique ID
func (gameserver *Gameserver) BeforeCreate(tx *gorm.DB) (err error) {
	if gameserver.ID == uuid.Nil {
		gameserver.ID = uuid.New()
	}
	var validate = validator.New()

	if err := validate.Struct(gameserver); err != nil {
		return err
	}

	return
}
