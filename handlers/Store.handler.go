package handlers

import (
	"0xKowalski1/server-hosting-web/config"
	"0xKowalski1/server-hosting-web/models"
	"0xKowalski1/server-hosting-web/services"
	"0xKowalski1/server-hosting-web/templates"
	"fmt"

	"log"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

type StoreHandler struct {
	CurrencyService *services.CurrencyService
	PriceService    *services.PriceService
}

func NewStoreHandler(currencyService *services.CurrencyService, priceService *services.PriceService) *StoreHandler {
	return &StoreHandler{
		CurrencyService: currencyService,
		PriceService:    priceService,
	}
}

func (sh *StoreHandler) GetStore(c echo.Context) error {
	priceMap, err := sh.getPrices(c)
	if err != nil {
		//500
	}

	return Render(c, 200, templates.StorePage(priceMap["memory"], priceMap["storage"], priceMap["archive"]))
}

func (sh *StoreHandler) SubmitStoreForm(c echo.Context) error {
	memory, memErr := strconv.Atoi(c.FormValue("memory"))
	storage, stoErr := strconv.Atoi(c.FormValue("storage"))
	archive, arcErr := strconv.Atoi(c.FormValue("archive"))

	prices, err := sh.getPrices(c)
	if err != nil {
		//500
	}

	// Validate
	if memErr != nil || stoErr != nil || arcErr != nil {
		// 400
	}

	// Init Stripe
	stripe.Key = config.Envs.StripeSecretKey

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

	stripeCheckoutSession, err := session.New(params)

	if err != nil {
		log.Printf("session.New: %v", err)
		// render error
	}

	c.Response().Header().Set("HX-Redirect", stripeCheckoutSession.URL)
	return c.NoContent(302)
}

func (sh *StoreHandler) GetGuidedStoreFlow(c echo.Context) error {
	return Render(c, 200, templates.GuidedStoreFlow())
}

func (sh *StoreHandler) GetAdvancedStoreFlow(c echo.Context) error {
	priceMap, err := sh.getPrices(c)
	if err != nil {
		//500
	}

	return Render(c, 200, templates.AdvancedStoreFlow(priceMap["memory"], priceMap["storage"], priceMap["archive"]))
}

func (sh *StoreHandler) getPrices(c echo.Context) (map[string]models.Price, error) {
	var user *models.User
	userInterface := c.Get("user")
	if userInterface != nil {
		userConversion, ok := userInterface.(*models.User)
		if ok {
			user = userConversion
		}
	}

	var currency models.Currency
	var err error // Avoid variable shadowing
	if user != nil {
		currency, err = sh.CurrencyService.GetCurrencyById(user.CurrencyID)

		if err != nil {
			return nil, err
		}
	} else {
		// Default to USD
		currency, err = sh.CurrencyService.GetDefaultCurrency()
	}

	// Get prices for currency
	prices, err := sh.PriceService.GetPrices(currency)
	if err != nil {
		return nil, err
	}
	priceMap := make(map[string]models.Price)
	for _, price := range prices {
		priceMap[price.Type] = price
	}

	return priceMap, nil
}
