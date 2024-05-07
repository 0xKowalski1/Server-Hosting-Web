package services

import (
	"0xKowalski1/server-hosting-web/models"
	"fmt"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (service *UserService) FindOrCreateUser(unknownUser models.User) (*models.User, error) {
	var user *models.User
	result := service.DB.Where("email = ? AND provider = ?", unknownUser.Email, unknownUser.Provider).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Set currency to USD by default, REMOVE ME, this should be handled elsewhere
			var usd = models.Currency{}
			service.DB.First(&usd, "code = ?", "USD")

			unknownUser.CurrencyID = usd.ID

			return service.CreateUser(unknownUser)
		}

		return nil, fmt.Errorf("Internal server error when trying to find user during FindOrCreateUser: %v", result.Error)
	}

	return user, nil
}

func (service *UserService) GetUser(userID string) (*models.User, error) {
	var user *models.User
	result := service.DB.Preload("Subscription").Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("User not found with ID: %s, due to error: %v", userID, result.Error)
	}
	return user, nil
}

func (service *UserService) CreateUser(newUser models.User) (*models.User, error) {
	user := &newUser
	result := service.DB.Create(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("Failed to create new user: %v", result.Error)
	}
	return user, nil
}
