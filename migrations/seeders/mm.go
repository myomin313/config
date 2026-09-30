package seeders

import (
	"log"

	"config-service/config"
	"config-service/models"

	"gorm.io/gorm"
)

func SeedMerchantCurrency() {

	defaultMerchantCurrencies := []models.MerchantCurrency{
		{
			ID:           1,
			MerchantCode: "M34794",
			Currency:     "USD",
			PlatformID:   1,
		},
		{
			ID:           4,
			MerchantCode: "M58697",
			Currency:     "THB",
			PlatformID:   1,
		},
		{
			ID:           5,
			MerchantCode: "M54452",
			Currency:     "USD",
			PlatformID:   1,
		},
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {

		// Delete all existing records
		for _, merchantCurrency := range defaultMerchantCurrencies {

			var existing models.MerchantCurrency

			err := tx.
				Where(
					"platform_id = ? AND merchant_code = ? AND currency = ?",
					merchantCurrency.PlatformID,
					merchantCurrency.MerchantCode,
					merchantCurrency.Currency,
				).
				First(&existing).Error

			if err == nil {
				// already exists, skip
				continue
			}

			if err != gorm.ErrRecordNotFound {
				return err
			}

			if err := tx.Create(&merchantCurrency).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("❌ Failed to seed MerchantCurrency: %v\n", err)
		return
	}

	log.Println("✅ MerchantCurrency seeded successfully")
}
