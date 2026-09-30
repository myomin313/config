package seeders

import (
	"log"

	"config-service/config"
	"config-service/models"

	"gorm.io/gorm"
)

func SeedGameData() {

	defaultGameData := []models.GameData{
		{
			ID:     101,
			NameCh: "百家乐",
			NameEn: "BACCARAT",
		},
		{
			ID:     102,
			NameCh: "龙虎",
			NameEn: "DRAGON TIGER",
		},
		{
			ID:     103,
			NameCh: "牛牛",
			NameEn: "NIU NIU",
		},
		{
			ID:     105,
			NameCh: "鱼虾蟹",
			NameEn: "KLA KLOUK",
		},
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {

		// 🔥 Step 1: Delete ALL existing data
		if err := tx.Where("1 = 1").
			Delete(&models.GameData{}).Error; err != nil {
			return err
		}

		// 🔥 Step 2: Insert fresh data
		if err := tx.Create(&defaultGameData).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("❌ Failed to seed GameData: %v\n", err)
		return
	}

	log.Println("✅ GameData seeded successfully")
}
