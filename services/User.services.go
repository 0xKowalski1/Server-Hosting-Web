package services

import (
	"0xKowalski1/server-hosting-web/models"
	"log"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (service *UserService) FindOrCreateUser(unknownUser models.User) (models.User, error) {
	var user models.User
	// Attempt to find the user first by email and provider
	result := service.DB.Where("email = ? AND provider = ?", unknownUser.Email, unknownUser.Provider).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// User not found, create new user

			// Set currency to USD by default, REMOVE ME, this should be handled elsewhere
			var usd = models.Currency{}
			service.DB.First(&usd, "code = ?", "USD")

			unknownUser.CurrencyID = usd.ID

			return service.CreateUser(unknownUser)
		}
		// Some other error occurred
		return models.User{}, result.Error
	}
	// User found, return existing user
	return user, nil
}

func (service *UserService) GetUser(userID string) (*models.User, error) {
	var user *models.User
	result := service.DB.Preload("Subscription").Where("id = ?", userID).First(&user)
	if result.Error != nil {
		log.Printf("User not found: %s, error: %v", userID, result.Error)
		return nil, result.Error
	}
	return user, nil
}

func (service *UserService) CreateUser(newUser models.User) (models.User, error) {
	result := service.DB.Create(&newUser) // Use GORM's Create method to add a new record
	if result.Error != nil {
		return models.User{}, result.Error // Return an empty User struct and the error
	}
	return newUser, nil // Return the newly created user and nil for the error
}
