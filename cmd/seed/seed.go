package main

import (
	"0xKowalski1/server-hosting-web/models"
	"0xKowalski1/server-hosting-web/seeds"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=password dbname=hosting port=5432 sslmode=disable TimeZone=Europe/London"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Seed Game
	db.Migrator().DropTable(&models.Game{})
	db.AutoMigrate(&models.Game{})
	games := seeds.SeedGames()
	for _, game := range games {
		db.Create(&game)
	}
}
