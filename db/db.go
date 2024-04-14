package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// InitDB initializes and returns a GORM database connection
func InitDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=password dbname=hosting port=5432 sslmode=disable TimeZone=Europe/London"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	return db
}
