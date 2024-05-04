package seeds

import "0xKowalski1/server-hosting-web/models"

func SeedCurrency() []*models.Currency {
	var seeds []*models.Currency

	usd := &models.Currency{
		Code:   "USD",
		Symbol: "$",
	}

	seeds = append(seeds, usd)

	gbp := &models.Currency{
		Code:   "GBP",
		Symbol: "£",
	}

	seeds = append(seeds, gbp)

	return seeds
}
