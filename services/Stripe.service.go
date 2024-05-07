package services

import (
	"0xKowalski1/server-hosting-web/models"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/subscription"
	"gorm.io/gorm"
)

type StripeService struct {
	DB *gorm.DB
}

func NewStripeService(db *gorm.DB) *StripeService {
	return &StripeService{
		DB: db,
	}
}

func (service *StripeService) CreateCheckoutSession(memory, storage, archive int, prices map[string]models.Price) (*stripe.CheckoutSession, error) {
	memoryCost := memory * prices["memory"].PricePerUnit
	storageCost := storage * prices["storage"].PricePerUnit
	archiveCost := archive * prices["archive"].PricePerUnit
	totalPrice := memoryCost + storageCost + archiveCost

	metadata := map[string]string{
		"memory_gb":       strconv.Itoa(memory),
		"memory_price_id": prices["memory"].ID.String(),

		"storage_gb":       strconv.Itoa(storage),
		"storage_price_id": prices["storage"].ID.String(),

		"archive_gb":       strconv.Itoa(archive),
		"archive_price_id": prices["archive"].ID.String(),
	}

	description := fmt.Sprintf(
		"Your Gameserver package consists of: %d GB of Memory at $%.2f a month, %d GB of Storage Space at $%.2f a month and %d GB of Archive Space at $%.2f a month.",
		memory, float64(memoryCost)/100, storage, float64(storageCost)/100, archive, float64(archiveCost)/100)

	domain := "http://localhost:3000"
	params := &stripe.CheckoutSessionParams{
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: metadata,
		},
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:        stripe.String("InterstellarHosts Gamesevers"),
						Description: &description,
					},
					Recurring: &stripe.CheckoutSessionLineItemPriceDataRecurringParams{
						Interval:      stripe.String("month"),
						IntervalCount: stripe.Int64(1),
					},
					UnitAmount: stripe.Int64(int64(totalPrice)),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(domain + "/store/callback"),
		CancelURL:  stripe.String(domain + "/store"),
	}

	session, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("Failed to create a new checkout session: %v", err)
	}

	return session, nil
}

func (service *StripeService) CreateSubscription(stripeSubscription *stripe.Subscription, user models.User) (*models.Subscription, error) {
	memoryGB, errMemory := strconv.Atoi(stripeSubscription.Metadata["memory_gb"])
	if errMemory != nil {
		return nil, fmt.Errorf("error converting memory_gb to int: %v", errMemory)
	}

	memoryPriceID, errMemoryPrice := uuid.Parse(stripeSubscription.Metadata["memory_price_id"])
	if errMemoryPrice != nil {
		return nil, fmt.Errorf("error parsing memory_price_id UUID: %v", errMemoryPrice)
	}

	storageGB, errStorage := strconv.Atoi(stripeSubscription.Metadata["storage_gb"])
	if errStorage != nil {
		return nil, fmt.Errorf("error converting storage_gb to int: %v", errStorage)
	}

	storagePriceID, errStoragePrice := uuid.Parse(stripeSubscription.Metadata["storage_price_id"])
	if errStoragePrice != nil {
		return nil, fmt.Errorf("error parsing storage_price_id UUID: %v", errStoragePrice)
	}

	archiveGB, errArchive := strconv.Atoi(stripeSubscription.Metadata["archive_gb"])
	if errArchive != nil {
		return nil, fmt.Errorf("error converting archive_gb to int: %v", errArchive)
	}

	archivePriceID, errArchivePrice := uuid.Parse(stripeSubscription.Metadata["archive_price_id"])
	if errArchivePrice != nil {
		return nil, fmt.Errorf("error parsing archive_price_id UUID: %v", errArchivePrice)
	}

	subscription := &models.Subscription{
		ID: stripeSubscription.ID,

		UserID: user.ID,

		Status: stripeSubscription.Status,

		MemoryGB:      memoryGB,
		MemoryPriceID: memoryPriceID,

		StorageGB:      storageGB,
		StoragePriceID: storagePriceID,

		ArchiveGB:      archiveGB,
		ArchivePriceID: archivePriceID,
	}

	result := service.DB.Create(&subscription)
	if result.Error != nil {
		return nil, fmt.Errorf("error creating subscription in the database: %v", result.Error)
	}

	return subscription, nil
}

func (service *StripeService) GetStripeSession(sessionID string) (*stripe.CheckoutSession, error) {
	session, err := session.Get(sessionID, nil)

	if err != nil {
		return nil, fmt.Errorf("Failed to get stripe session with ID: %s due to error: %v", sessionID, err)
	}

	return session, nil
}

func (service *StripeService) GetStripeSubscription(subscriptionID string) (*stripe.Subscription, error) {

	subscription, err := subscription.Get(subscriptionID, nil)

	if err != nil {
		return nil, fmt.Errorf("Failed to get stripe subscription at ID: %s due to error: %v", subscriptionID, err)
	}

	return subscription, nil
}

func (service *StripeService) GetSubscriptionByUserID(userID string) (*models.Subscription, error) {
	var subscription *models.Subscription

	result := service.DB.First(&subscription, "user_id = ?", userID)
	if result.Error != nil {
		return nil, fmt.Errorf("Failed to query database for subscription at UserID: %s due to error: %v", userID, result.Error)
	}
	return subscription, nil
}

// Redundant to return a bool and an error, so error = Cannot allocate
func (service *StripeService) CanAllocateResources(userID string, memoryGB, storageGB int, gameserverService *GameserverService) error {
	subscription, err := service.GetSubscriptionByUserID(userID)
	if err != nil {
		return fmt.Errorf("Cannot allocate resources as we cant find a subscription with UserID: %s due to error: %v", userID, err)
	}

	gameservers, err := gameserverService.GetGameservers(userID)
	if err != nil {
		return fmt.Errorf("Cannot allocate resources as query for gameservers with UserID: %s failed due to error: %v", userID, err)
	}

	usedMemory, usedStorage := gameserverService.GetUsedResources(gameservers)

	if ((usedMemory + memoryGB) > subscription.MemoryGB) || ((usedStorage + storageGB) > subscription.StorageGB) {
		return fmt.Errorf("Cannot allocate resources as requested resources exceeds subscription limits")
	}

	return nil
}
