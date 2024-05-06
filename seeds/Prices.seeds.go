package seeds

import "0xKowalski1/server-hosting-web/models"

func SeedPrice(currencies []*models.Currency) []models.Price {
	var seeds []models.Price
	currencyMap := make(map[string]*models.Currency)
	for _, currency := range currencies {
		currencyMap[currency.Code] = currency
	}

	memoryUSD := models.Price{
		Type:         "memory",
		PricePerUnit: 500,
		CurrencyID:   currencyMap["USD"].ID,
		StripeID:     "price_1NpdWeAgE84zSyWHqkZokTfL",
	}
	storageUSD := models.Price{
		Type:         "storage",
		PricePerUnit: 50,
		CurrencyID:   currencyMap["USD"].ID,
		StripeID:     "price_1NpdWeAgE84zSyWH6DYtNKO8",
	}
	archiveUSD := models.Price{
		Type:         "archive",
		PricePerUnit: 10,
		CurrencyID:   currencyMap["USD"].ID,
		StripeID:     "price_1NpdWeAgE84zSyWHYuSFvMCl",
	}

	seeds = append(seeds, memoryUSD)
	seeds = append(seeds, storageUSD)
	seeds = append(seeds, archiveUSD)

	memoryGBP := models.Price{
		Type:         "memory",
		PricePerUnit: 400,
		CurrencyID:   currencyMap["GBP"].ID,
		StripeID:     "price_1NpdWeAgE84zSyWHsdqrb5Sv",
	}
	storageGBP := models.Price{
		Type:         "storage",
		PricePerUnit: 40,
		CurrencyID:   currencyMap["GBP"].ID,
		StripeID:     "price_1NpdWeAgE84zSyWHNNtxebnm",
	}
	archiveGBP := models.Price{
		Type:         "archive",
		PricePerUnit: 8,
		CurrencyID:   currencyMap["GBP"].ID,
		StripeID:     "price_1NpdWeAgE84zSyWHAWzr6PWG",
	}

	seeds = append(seeds, memoryGBP)
	seeds = append(seeds, storageGBP)
	seeds = append(seeds, archiveGBP)

	return seeds
}
