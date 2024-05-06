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
	}
	storageUSD := models.Price{
		Type:         "storage",
		PricePerUnit: 50,
		CurrencyID:   currencyMap["USD"].ID,
	}
	archiveUSD := models.Price{
		Type:         "archive",
		PricePerUnit: 10,
		CurrencyID:   currencyMap["USD"].ID,
	}

	seeds = append(seeds, memoryUSD)
	seeds = append(seeds, storageUSD)
	seeds = append(seeds, archiveUSD)

	memoryGBP := models.Price{
		Type:         "memory",
		PricePerUnit: 400,
		CurrencyID:   currencyMap["GBP"].ID,
	}
	storageGBP := models.Price{
		Type:         "storage",
		PricePerUnit: 40,
		CurrencyID:   currencyMap["GBP"].ID,
	}
	archiveGBP := models.Price{
		Type:         "archive",
		PricePerUnit: 8,
		CurrencyID:   currencyMap["GBP"].ID,
	}

	seeds = append(seeds, memoryGBP)
	seeds = append(seeds, storageGBP)
	seeds = append(seeds, archiveGBP)

	return seeds
}
