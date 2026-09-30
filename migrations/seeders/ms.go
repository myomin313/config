package seeders

import (
	"config-service/config"
	"config-service/models"
	"errors"
	"log"

	"gorm.io/gorm"
)

func CreateMasterGameTypes() {
	log.Println("=== Starting MasterGameTypes Migration ===")

	masterGameTypes := []models.MasterGameType{
		{ID: 1, GameTypeID: 1, Odds: 1, Image: "tiger.png", KeyName: "tiger"},
		{ID: 2, GameTypeID: 1, Odds: 1, Image: "klouk.png", KeyName: "klouk"},
		{ID: 3, GameTypeID: 1, Odds: 1, Image: "chicken.png", KeyName: "chicken"},
		{ID: 4, GameTypeID: 1, Odds: 1, Image: "shrimp.png", KeyName: "shrimp"},
		{ID: 5, GameTypeID: 1, Odds: 1, Image: "fish.png", KeyName: "fish"},
		{ID: 6, GameTypeID: 1, Odds: 1, Image: "crab.png", KeyName: "crab"},
		// {ID: 7, GameTypeID: 2, Odds: 1, Image: "triple_color.png", KeyName: "triple_color"},
		// {ID: 8, GameTypeID: 2, Odds: 1, Image: "triple_symbol.png", KeyName: "triple_symbol"},
		{ID: 9, GameTypeID: 1, Odds: 1, Image: "triple_color.png", KeyName: "triple_color"},
		{ID: 10, GameTypeID: 1, Odds: 1, Image: "triple_symbol.png", KeyName: "triple_symbol"},
		{ID: 11, GameTypeID: 2, Odds: 1, Image: "tiger.png", KeyName: "tiger"},
		{ID: 12, GameTypeID: 2, Odds: 1, Image: "klouk.png", KeyName: "klouk"},
		{ID: 13, GameTypeID: 2, Odds: 1, Image: "chicken.png", KeyName: "chicken"},
		{ID: 14, GameTypeID: 2, Odds: 1, Image: "shrimp.png", KeyName: "shrimp"},
		{ID: 15, GameTypeID: 2, Odds: 1, Image: "fish.png", KeyName: "fish"},
		{ID: 16, GameTypeID: 2, Odds: 1, Image: "crab.png", KeyName: "crab"},
		// {ID: 17, GameTypeID: 3, Odds: 1, Image: "klouk.png", KeyName: "klouk"},
		// {ID: 18, GameTypeID: 3, Odds: 1, Image: "tiger.png", KeyName: "tiger"},
		// {ID: 19, GameTypeID: 3, Odds: 1, Image: "chicken.png", KeyName: "chicken"},
		// {ID: 20, GameTypeID: 3, Odds: 1, Image: "shrimp.png", KeyName: "shrimp"},
		// {ID: 21, GameTypeID: 3, Odds: 1, Image: "fish.png", KeyName: "fish"},
		// {ID: 22, GameTypeID: 3, Odds: 1, Image: "crab.png", KeyName: "crab"},
		// {ID: 23, GameTypeID: 3, Odds: 1, Image: "triple_color.png", KeyName: "triple_color"},
		// {ID: 24, GameTypeID: 3, Odds: 1, Image: "triple_symbol.png", KeyName: "triple_symbol"},
		{ID: 30, GameTypeID: 5, Odds: 1, Image: "dragon_one_dragon.png", KeyName: "dragon_one_dragon"},
		{ID: 31, GameTypeID: 5, Odds: 1, Image: "dragon_one_tiger.png", KeyName: "dragon_one_tiger"},
		{ID: 32, GameTypeID: 5, Odds: 1, Image: "dragon_one_pair.png", KeyName: "dragon_one_pair"},
		{ID: 33, GameTypeID: 5, Odds: 1, Image: "dragon_one_flush_pair.png", KeyName: "dragon_one_flush_pair"},
		{ID: 34, GameTypeID: 6, Odds: 1, Image: "dragon_two_dragon.png", KeyName: "dragon_two_dragon"},
		{ID: 35, GameTypeID: 6, Odds: 1, Image: "dragon_two_tiger.png", KeyName: "dragon_two_tiger"},
		{ID: 36, GameTypeID: 6, Odds: 1, Image: "dragon_two_tie.png", KeyName: "dragon_two_tie"},
		{ID: 37, GameTypeID: 6, Odds: 1, Image: "double_animal.png", KeyName: "double_animal"},
		{ID: 38, GameTypeID: 5, Odds: 1, Image: "double_animal.png", KeyName: "double_animal"},

		{ID: 40, GameTypeID: 7, Odds: 1, Image: "gold_coin.png", KeyName: "gold_coin"},
		{ID: 41, GameTypeID: 7, Odds: 1, Image: "klouk.png", KeyName: "klouk"},
		{ID: 42, GameTypeID: 7, Odds: 1, Image: "chicken.png", KeyName: "chicken"},
		{ID: 43, GameTypeID: 7, Odds: 1, Image: "shrimp.png", KeyName: "shrimp"},
		{ID: 44, GameTypeID: 7, Odds: 1, Image: "fish.png", KeyName: "fish"},
		{ID: 45, GameTypeID: 7, Odds: 1, Image: "crab.png", KeyName: "crab"},
		{ID: 46, GameTypeID: 7, Odds: 1, Image: "triple_symbol.png", KeyName: "triple_symbol"},
		{ID: 47, GameTypeID: 7, Odds: 1, Image: "double_animal.png", KeyName: "double_animal"},
	}
	translations := map[uint64][]models.GameTypeTranslation{
		1: {
			{LanguageCode: "en", Name: "Tiger"},
			{LanguageCode: "ch", Name: "老虎"},
			{LanguageCode: "kh", Name: "ខ្លា"},
		},
		2: {
			{LanguageCode: "en", Name: "Gourd"},
			{LanguageCode: "ch", Name: "葫芦"},
			{LanguageCode: "kh", Name: "ឃ្លោក"},
		},
		3: {
			{LanguageCode: "en", Name: "Chicken"},
			{LanguageCode: "ch", Name: "鸡"},
			{LanguageCode: "kh", Name: "មាន់"},
		},
		4: {
			{LanguageCode: "en", Name: "Shrimp"},
			{LanguageCode: "ch", Name: "虾"},
			{LanguageCode: "kh", Name: "បង្គា"},
		},
		5: {
			{LanguageCode: "en", Name: "Fish"},
			{LanguageCode: "ch", Name: "鱼"},
			{LanguageCode: "kh", Name: "ត្រី"},
		},
		6: {
			{LanguageCode: "en", Name: "Crab"},
			{LanguageCode: "ch", Name: "蟹"},
			{LanguageCode: "kh", Name: "ក្តាម"},
		},
		7: {
			{LanguageCode: "en", Name: "Any Triple Colors"},
			{LanguageCode: "ch", Name: "任意三色"},
			{LanguageCode: "kh", Name: "ពណ៌បី"},
		},
		8: {
			{LanguageCode: "en", Name: "Any Triple Symbols"},
			{LanguageCode: "ch", Name: "豹子"},
			{LanguageCode: "kh", Name: "និមិត្តសញ្ញាបី"},
		},
		9: {
			{LanguageCode: "en", Name: "Any Triple Colors"},
			{LanguageCode: "ch", Name: "任意三色"},
			{LanguageCode: "kh", Name: "ពណ៌បី"},
		},
		10: {
			{LanguageCode: "en", Name: "Any Triple Symbols"},
			{LanguageCode: "ch", Name: "豹子"},
			{LanguageCode: "kh", Name: "និមិត្តសញ្ញាបី"},
		},
		11: {
			{LanguageCode: "en", Name: "Tiger"},
			{LanguageCode: "ch", Name: "老虎"},
			{LanguageCode: "kh", Name: "ខ្លា"},
		},
		12: {
			{LanguageCode: "en", Name: "Gourd"},
			{LanguageCode: "ch", Name: "葫芦"},
			{LanguageCode: "kh", Name: "ឃ្លោក"},
		},
		13: {
			{LanguageCode: "en", Name: "Chicken"},
			{LanguageCode: "ch", Name: "鸡"},
			{LanguageCode: "kh", Name: "មាន់"},
		},
		14: {
			{LanguageCode: "en", Name: "Shrimp"},
			{LanguageCode: "ch", Name: "虾"},
			{LanguageCode: "kh", Name: "បង្គា"},
		},
		15: {
			{LanguageCode: "en", Name: "Fish"},
			{LanguageCode: "ch", Name: "鱼"},
			{LanguageCode: "kh", Name: "ត្រី"},
		},
		16: {
			{LanguageCode: "en", Name: "Crab"},
			{LanguageCode: "ch", Name: "蟹"},
			{LanguageCode: "kh", Name: "ក្តាម"},
		},
		17: {
			{LanguageCode: "en", Name: "Gourd"},
			{LanguageCode: "ch", Name: "葫芦"},
			{LanguageCode: "kh", Name: "ឃ្លោក"},
		},
		18: {
			{LanguageCode: "en", Name: "Tiger"},
			{LanguageCode: "ch", Name: "老虎"},
			{LanguageCode: "kh", Name: "ខ្លា"},
		},
		19: {
			{LanguageCode: "en", Name: "Chicken"},
			{LanguageCode: "ch", Name: "鸡"},
			{LanguageCode: "kh", Name: "មាន់"},
		},
		20: {
			{LanguageCode: "en", Name: "Shrimp"},
			{LanguageCode: "ch", Name: "虾"},
			{LanguageCode: "kh", Name: "បង្គា"},
		},
		21: {
			{LanguageCode: "en", Name: "Fish"},
			{LanguageCode: "ch", Name: "鱼"},
			{LanguageCode: "kh", Name: "ត្រី"},
		},
		22: {
			{LanguageCode: "en", Name: "Crab"},
			{LanguageCode: "ch", Name: "蟹"},
			{LanguageCode: "kh", Name: "ក្តាម"},
		},
		23: {
			{LanguageCode: "en", Name: "Any Triple Colors"},
			{LanguageCode: "ch", Name: "任意三色"},
			{LanguageCode: "kh", Name: "ពណ៌បី"},
		},
		24: {
			{LanguageCode: "en", Name: "Any Triple Symbols"},
			{LanguageCode: "ch", Name: "豹子"},
			{LanguageCode: "kh", Name: "និមិត្តសញ្ញាបី"},
		},
		30: {
			{LanguageCode: "en", Name: "DRAGON"},
			{LanguageCode: "ch", Name: "龙"},
			{LanguageCode: "kh", Name: "ខ្លា"},
		},
		31: {
			{LanguageCode: "en", Name: "TIGER"},
			{LanguageCode: "ch", Name: "虎"},
			{LanguageCode: "kh", Name: "មាន់"},
		},
		32: {
			{LanguageCode: "en", Name: "PAIR"},
			{LanguageCode: "ch", Name: "对子"},
			{LanguageCode: "kh", Name: "ក្រពើ"},
		},
		33: {
			{LanguageCode: "en", Name: "FLUSH PAIR"},
			{LanguageCode: "ch", Name: "同花对子"},
			{LanguageCode: "kh", Name: "ត្រី"},
		},
		34: {
			{LanguageCode: "en", Name: "DRAGON"},
			{LanguageCode: "ch", Name: "龙"},
			{LanguageCode: "kh", Name: "មឹក"},
		},
		35: {
			{LanguageCode: "en", Name: "TIGER"},
			{LanguageCode: "ch", Name: "虎"},
			{LanguageCode: "kh", Name: "ពណ៌បី"},
		},
		36: {
			{LanguageCode: "en", Name: "TIE"},
			{LanguageCode: "ch", Name: "和"},
			{LanguageCode: "kh", Name: "និមិត្តសញ្ញាបី"},
		},
		37: {
			{LanguageCode: "en", Name: "Double Dices"},
			{LanguageCode: "ch", Name: "双骰"},
			{LanguageCode: "kh", Name: "គ្រាប់ឡុកឡាក់ពីរ"},
		},
		38: {
			{LanguageCode: "en", Name: "Double Dices"},
			{LanguageCode: "ch", Name: "双骰"},
			{LanguageCode: "kh", Name: "គ្រាប់ឡុកឡាក់ពីរ"},
		},

		40: {
			{LanguageCode: "en", Name: "Gold Coin"},
			{LanguageCode: "ch", Name: "金币"},
			{LanguageCode: "kh", Name: "កាក់មាស"},
		},
		41: {
			{LanguageCode: "en", Name: "Gourd"},
			{LanguageCode: "ch", Name: "葫芦"},
			{LanguageCode: "kh", Name: "ឃ្លោក"},
		},
		42: {
			{LanguageCode: "en", Name: "Chicken"},
			{LanguageCode: "ch", Name: "鸡"},
			{LanguageCode: "kh", Name: "មាន់"},
		},
		43: {
			{LanguageCode: "en", Name: "Shrimp"},
			{LanguageCode: "ch", Name: "虾"},
			{LanguageCode: "kh", Name: "បង្គា"},
		},
		44: {
			{LanguageCode: "en", Name: "Fish"},
			{LanguageCode: "ch", Name: "鱼"},
			{LanguageCode: "kh", Name: "ត្រី"},
		},
		45: {
			{LanguageCode: "en", Name: "Crab"},
			{LanguageCode: "ch", Name: "蟹"},
			{LanguageCode: "kh", Name: "ក្តាម"},
		},
		46: {
			{LanguageCode: "en", Name: "Any Triple Symbols"},
			{LanguageCode: "ch", Name: "豹子"},
			{LanguageCode: "kh", Name: "និមិត្តសញ្ញាបី"},
		},
		47: {
			{LanguageCode: "en", Name: "Double Dices"},
			{LanguageCode: "ch", Name: "双骰"},
			{LanguageCode: "kh", Name: "គ្រាប់ឡុកឡាក់ពីរ"},
		},
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		for _, gt := range masterGameTypes {
			log.Printf("Processing MasterGameType ID: %d | KeyName: %s", gt.ID, gt.KeyName)

			var existing models.MasterGameType
			err := tx.Where("id = ?", gt.ID).First(&existing).Error

			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				if err := tx.Create(&gt).Error; err != nil {
					return err
				}
				log.Printf("Created MasterGameType ID: %d", gt.ID)

			case err != nil:
				return err

			default:
				if err := tx.Model(&existing).Updates(map[string]interface{}{
					"game_type_id": gt.GameTypeID,
					"odds":         gt.Odds,
					"image":        gt.Image,
					"key_name":     gt.KeyName,
				}).Error; err != nil {
					return err
				}
				log.Printf("Updated MasterGameType ID: %d", gt.ID)
			}

			if trs, ok := translations[gt.ID]; ok {
				for _, tr := range trs {
					tr.MasterGameTypeID = gt.ID

					var existingTr models.GameTypeTranslation
					err := tx.
						Where("master_game_type_id = ? AND language_code = ?", gt.ID, tr.LanguageCode).
						First(&existingTr).Error

					switch {
					case errors.Is(err, gorm.ErrRecordNotFound):
						if err := tx.Create(&tr).Error; err != nil {
							return err
						}
						log.Printf("Created translation | MasterGameTypeID: %d | Lang: %s", gt.ID, tr.LanguageCode)

					case err != nil:
						return err

					default:
						if err := tx.Model(&existingTr).Updates(map[string]interface{}{
							"name": tr.Name,
						}).Error; err != nil {
							return err
						}
						log.Printf("Updated translation | MasterGameTypeID: %d | Lang: %s", gt.ID, tr.LanguageCode)
					}
				}
			} else {
				log.Printf("No translations found for ID: %d", gt.ID)
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("=== MasterGameTypes Migration Failed: %v ===", err)
		return
	}

	log.Println("=== MasterGameTypes Migration Completed ===")
}
