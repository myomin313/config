package seeders

// import (
// 	"log"

// 	"config-service/config"
// 	"config-service/models"

// 	"gorm.io/gorm"
// )

// func SeedCurrency() {

// 	defaultCurrencies := []models.Currency{
// 		{ID: 1, Code: "USD", Name: "US Dollar", Symbol: "$", Precision: 2, Status: "enabled"},
// 		{ID: 2, Code: "KHR", Name: "Cambodian Riel", Symbol: "៛", Precision: 0, Status: "enabled"},
// 		{ID: 3, Code: "CNY", Name: "Chinese Yuan", Symbol: "¥", Precision: 2, Status: "enabled"},
// 		{ID: 4, Code: "THB", Name: "Thai Baht", Symbol: "฿", Precision: 2, Status: "enabled"},
// 		{ID: 5, Code: "VND", Name: "Vietnamese Dong", Symbol: "₫", Precision: 0, Status: "enabled"},
// 		{ID: 6, Code: "SGD", Name: "Singapore Dollar", Symbol: "S$", Precision: 2, Status: "enabled"},
// 		{ID: 7, Code: "MYR", Name: "Malaysian Ringgit", Symbol: "RM", Precision: 2, Status: "enabled"},
// 	}

// 	err := config.DB.Transaction(func(tx *gorm.DB) error {

// 		for _, currency := range defaultCurrencies {

// 			var existing models.Currency

// 			err := tx.
// 				Where("code = ?", currency.Code).
// 				First(&existing).Error

// 			// already exists → skip
// 			if err == nil {
// 				continue
// 			}

// 			// real error
// 			if err != gorm.ErrRecordNotFound {
// 				return err
// 			}

// 			// insert new
// 			if err := tx.Create(&currency).Error; err != nil {
// 				return err
// 			}
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		log.Printf("❌ Failed to seed Currency: %v\n", err)
// 		return
// 	}

// 	log.Println("✅ Currency seeded successfully")
// }
