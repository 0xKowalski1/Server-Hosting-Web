package services

import (
	"0xKowalski1/server-hosting-web/config"
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
	// Init Stripe
	stripe.Key = config.Envs.StripeSecretKey

	return &StripeService{
		DB: db,
	}
}

func (ss *StripeService) CreateCheckoutSession(memory, storage, archive int, prices map[string]models.Price) (*stripe.CheckoutSession, error) {
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

	return session.New(params)
}

func (ss *StripeService) CreateSubscription(stripeSubscription *stripe.Subscription, user models.User) (*models.Subscription, error) {
	// Pain
	memoryGB, errMemory := strconv.Atoi(stripeSubscription.Metadata["memory_gb"])
	if errMemory != nil {
		return nil, errMemory
	}
	memoryPriceID, errMemoryPrice := uuid.Parse(stripeSubscription.Metadata["memory_price_id"])
	if errMemoryPrice != nil {
		return nil, errMemoryPrice
	}

	storageGB, errStorage := strconv.Atoi(stripeSubscription.Metadata["storage_gb"])
	if errStorage != nil {
		return nil, errStorage
	}
	storagePriceID, errStoragePrice := uuid.Parse(stripeSubscription.Metadata["storage_price_id"])
	if errStoragePrice != nil {
		return nil, errStoragePrice
	}

	archiveGB, errArchive := strconv.Atoi(stripeSubscription.Metadata["archive_gb"])
	if errArchive != nil {
		return nil, errArchive
	}
	archivePriceID, errArchivePrice := uuid.Parse(stripeSubscription.Metadata["archive_price_id"])
	if errArchivePrice != nil {
		return nil, errArchivePrice
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

	result := ss.DB.Create(&subscription)
	if result.Error != nil {
		return nil, result.Error
	}

	return subscription, nil
}

func (ss *StripeService) GetStripeSession(sessionID string) (*stripe.CheckoutSession, error) {
	return session.Get(sessionID, nil)
}

func (ss *StripeService) GetStripeSubscription(subscriptionID string) (*stripe.Subscription, error) {
	return subscription.Get(subscriptionID, nil)
}
