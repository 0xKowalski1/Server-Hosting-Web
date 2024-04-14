package models

import (
	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	Name             string
	ShortDescription string
	GridImage        string
}
