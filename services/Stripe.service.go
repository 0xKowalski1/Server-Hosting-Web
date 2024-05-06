package services

import (
	"0xKowalski1/server-hosting-web/config"
	"0xKowalski1/server-hosting-web/models"
	"fmt"
	"strconv"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

type StripeService struct {
}

func NewStripeService() *StripeService {
	// Init Stripe
	stripe.Key = config.Envs.StripeSecretKey

	return &StripeService{}
}

func (ss *StripeService) CreateCheckoutSession(memory, storage, archive int, prices map[string]models.Price) (*stripe.CheckoutSession, error) {
	memoryCost := memory * prices["memory"].PricePerUnit
	storageCost := storage * prices["storage"].PricePerUnit
	archiveCost := archive * prices["archive"].PricePerUnit
	totalPrice := memoryCost + storageCost + archiveCost

	metadata := map[string]string{
		"memory_gb":           strconv.Itoa(memory),                        // Memory in GB
		"memory_cost":         strconv.Itoa(memoryCost),                    // Total cost for memory
		"memory_price_per_gb": strconv.Itoa(prices["memory"].PricePerUnit), // Cost per GB for memory

		"storage_gb":           strconv.Itoa(storage),                        // Storage in GB
		"storage_cost":         strconv.Itoa(storageCost),                    // Total cost for storage
		"storage_price_per_gb": strconv.Itoa(prices["storage"].PricePerUnit), // Cost per GB for storage

		"archive_gb":           strconv.Itoa(archive),                        // Archive space in GB
		"archive_cost":         strconv.Itoa(archiveCost),                    // Total cost for archive
		"archive_price_per_gb": strconv.Itoa(prices["archive"].PricePerUnit), // Cost per GB for archive
	}

	description := fmt.Sprintf(
		"Your Gameserver package consists of: %d GB of Memory at $%.2f a month, %d GB of Storage Space at $%.2f a month and %d GB of Archive Space at $%.2f a month.",
		memory, float64(memoryCost)/100, storage, float64(storageCost)/100, archive, float64(archiveCost)/100)

	domain := "http://localhost:3000"
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:        stripe.String("InterstellarHosts Gamesevers"),
						Description: &description,
						Metadata:    metadata,
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
		SuccessURL: stripe.String(domain + "/profile/gameservers"),
		CancelURL:  stripe.String(domain + "/store"),
	}

	return session.New(params)
}
