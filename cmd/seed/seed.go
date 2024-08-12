package main

import (
	"0xKowalski1/server-hosting-web/models"
	"0xKowalski1/server-hosting-web/seeds"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=password dbname=serverhosting port=5432 sslmode=disable TimeZone=Europe/London"
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

	// Seed User
	db.Migrator().DropTable(&models.User{})
	db.AutoMigrate(&models.User{})

	// Seed Gameserver
	db.Migrator().DropTable(&models.Gameserver{})
	db.AutoMigrate(&models.Gameserver{})

	// Seed Currency
	db.Migrator().DropTable(&models.Currency{})
	db.AutoMigrate(&models.Currency{})
	currencies := seeds.SeedCurrency()
	for _, currency := range currencies {
		db.Create(&currency)
	}

	// Seed Prices
	db.Migrator().DropTable(&models.Price{})
	db.AutoMigrate(&models.Price{})
	prices := seeds.SeedPrice(currencies)
	for _, price := range prices {
		db.Create(&price)
	}

	// Seed Subscriptions
	db.Migrator().DropTable(&models.Subscription{})
	db.AutoMigrate(&models.Subscription{})

}
