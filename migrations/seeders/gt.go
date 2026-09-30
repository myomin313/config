package seeders

import (
	"log"

	"config-service/config"
	"config-service/models"

	"gorm.io/gorm/clause"
)

func SeedGameTypes() {
	defaultGameTypes := []models.GameType{
		{ID: 1, GameTypeName: "Game Play Type 1", GameDataID: 105},
		{ID: 2, GameTypeName: "Game Play Type 2", GameDataID: 105},
		{ID: 5, GameTypeName: "Game Play Type 1", GameDataID: 102},
		{ID: 6, GameTypeName: "Game Play Type 2", GameDataID: 102},
		{ID: 7, GameTypeName: "Game Play Type 3", GameDataID: 105},
		//{ID: 3, GameTypeName: "Game Play Type 3",GameDataID:204},
	}

	for _, gt := range defaultGameTypes {
		err := config.DB.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}}, // conflict on primary key
			//DoUpdates: clause.AssignmentColumns([]string{"game_type_name"}),
			DoUpdates: clause.AssignmentColumns([]string{"game_type_name", "game_data_id"}),
		}).Create(&gt).Error

		if err != nil {
			log.Printf("Failed to upsert GameType ID %d: %v\n", gt.ID, err)
			continue
		}

		log.Printf("GameType ID %d seeded successfully\n", gt.ID)
	}
}
