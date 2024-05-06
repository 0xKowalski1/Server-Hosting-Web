package handlers

import (
	"0xKowalski1/server-hosting-web/models"
	"0xKowalski1/server-hosting-web/services"
	"0xKowalski1/server-hosting-web/templates"
	"net/http"
	"time"

	"log"
	"strconv"

	"github.com/labstack/echo/v4"
)

type StoreHandler struct {
	CurrencyService *services.CurrencyService
	PriceService    *services.PriceService
	StripeService   *services.StripeService
}

func NewStoreHandler(currencyService *services.CurrencyService, priceService *services.PriceService, stripeService *services.StripeService) *StoreHandler {
	return &StoreHandler{
		CurrencyService: currencyService,
		PriceService:    priceService,
		StripeService:   stripeService,
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

	stripeCheckoutSession, err := sh.StripeService.CreateCheckoutSession(memory, storage, archive, prices)
	if err != nil {
		log.Printf("session.New: %v", err)
		// render error
	}

	c.SetCookie(&http.Cookie{
		Name:     "stripeSessionID",
		Value:    stripeCheckoutSession.ID,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

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

func (sh *StoreHandler) StripeSuccessCallback(c echo.Context) error {
	cookie, err := c.Cookie("stripeSessionID")
	if err != nil || cookie.Value == "" {
		// Bad error!
		log.Println(err)

	}

	session, err := sh.StripeService.GetStripeSession(cookie.Value)
	if err != nil {
		// Bad error!
		log.Println(err)

	}

	subscription, err := sh.StripeService.GetStripeSubscription(session.Subscription.ID)
	if err != nil {
		log.Println(err)

	}

	var user *models.User
	userInterface := c.Get("user")
	if userInterface != nil {
		userConversion, ok := userInterface.(*models.User)
		if ok {
			user = userConversion
		}
	}
	if user == nil {
		// Bad error!
		log.Println(err)

	}

	_, err = sh.StripeService.CreateSubscription(subscription, *user)
	if err != nil {
		// Bad error!
		log.Println(err)
	}

	// Reset cookie
	c.SetCookie(&http.Cookie{
		Name:     "stripeSessionID",
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		Secure:   true, // Set to true in production
		SameSite: http.SameSiteLaxMode,
	})

	return c.Redirect(302, "/profile/gameservers")
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
	} else {
		// Default to USD
		currency, err = sh.CurrencyService.GetDefaultCurrency()
	}
	if err != nil {
		return nil, err
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
