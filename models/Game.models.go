package models

import (
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	ID               uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name             string
	ShortDescription string
	GridImage        string
	IconImage        string
	ContainerImage   string
}

// Set unique ID
func (game *Game) BeforeCreate(tx *gorm.DB) (err error) {
	if game.ID == uuid.Nil {
		game.ID = uuid.New()
	}
	return
}
