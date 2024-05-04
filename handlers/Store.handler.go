package handlers

import (
	"0xKowalski1/server-hosting-web/models"
	"0xKowalski1/server-hosting-web/services"
	"0xKowalski1/server-hosting-web/templates"
	"log"

	"github.com/labstack/echo/v4"
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

	log.Println(priceMap["memory"].Type)

	return Render(c, 200, templates.StorePage(priceMap["memory"], priceMap["storage"], priceMap["archive"]))
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

func (sh *StoreHandler) GetCheckout(c echo.Context) error {
	return Render(c, 200, templates.Checkout())
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
