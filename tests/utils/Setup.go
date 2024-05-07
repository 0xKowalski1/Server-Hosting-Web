package testutils

import (
	"0xKowalski1/server-hosting-web/models"

	"github.com/stripe/stripe-go/v78"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.Migrator().DropTable(&models.Game{})
	db.AutoMigrate(&models.Game{})

	db.Migrator().DropTable(&models.User{})
	db.AutoMigrate(&models.User{})

	db.Migrator().DropTable(&models.Gameserver{})
	db.AutoMigrate(&models.Gameserver{})

	db.Migrator().DropTable(&models.Currency{})
	db.AutoMigrate(&models.Currency{})

	db.Migrator().DropTable(&models.Price{})
	db.AutoMigrate(&models.Price{})

	db.Migrator().DropTable(&models.Subscription{})
	db.AutoMigrate(&models.Subscription{})

	if err != nil {
		panic(err)
	}

	return db
}

func SetupCurrency(db *gorm.DB) *models.Currency {
	currency := &models.Currency{
		Code:   "USD",
		Symbol: "$",
	}

	result := db.Create(currency)
	if result.Error != nil {
		panic(result.Error)
	}

	return currency
}

func SetupPrices(db *gorm.DB, currency *models.Currency) (*models.Price, *models.Price, *models.Price) {
	memoryUSD := &models.Price{
		Type:         "memory",
		PricePerUnit: 500,
		CurrencyID:   currency.ID,
	}
	storageUSD := &models.Price{
		Type:         "storage",
		PricePerUnit: 50,
		CurrencyID:   currency.ID,
	}
	archiveUSD := &models.Price{
		Type:         "archive",
		PricePerUnit: 10,
		CurrencyID:   currency.ID,
	}

	result := db.Create(memoryUSD)
	if result.Error != nil {
		panic(result.Error)
	}
	result = db.Create(storageUSD)
	if result.Error != nil {
		panic(result.Error)
	}
	result = db.Create(archiveUSD)
	if result.Error != nil {
		panic(result.Error)
	}

	return memoryUSD, storageUSD, archiveUSD
}

func SetupGame(db *gorm.DB) *models.Game {
	game := &models.Game{
		Name:             "Test Game",
		ShortDescription: "Test game description.",
		GridImage:        "/images/test-grid.jpg",
		IconImage:        "/images/test-icon.png",
		ContainerImage:   "testContainerImage",
	}

	result := db.Create(game)
	if result.Error != nil {
		panic(result.Error)
	}

	return game
}

func SetupUser(db *gorm.DB, currency *models.Currency) *models.User {
	user := &models.User{
		ID:         "1",
		Email:      "test@test.com",
		Provider:   "google",
		CurrencyID: currency.ID,
	}

	result := db.Create(user)
	if result.Error != nil {
		panic(result.Error)
	}

	return user
}

func SetupSubscription(db *gorm.DB, user *models.User, memoryPrice *models.Price, storagePrice *models.Price, archivePrice *models.Price) *models.Subscription {
	subscription := &models.Subscription{
		ID:        "1",
		UserID:    user.ID,
		Status:    stripe.SubscriptionStatusActive,
		MemoryGB:  10,
		StorageGB: 10,
		ArchiveGB: 10,

		MemoryPriceID:  memoryPrice.ID,
		StoragePriceID: storagePrice.ID,
		ArchivePriceID: archivePrice.ID,
	}

	result := db.Create(subscription)
	if result.Error != nil {
		panic(result.Error)
	}

	return subscription
}
