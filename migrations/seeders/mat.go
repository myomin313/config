package seeders

import (
	"config-service/config"
	"config-service/models"
	"log"
)

func Maintenance() error {
	adminID := uint(1)
	maintenance := models.Maintenance{
		ID:                     1,
		ApplicationMaintenance: false,
		CreatedByID:            &adminID,
		UpdatedByID:            &adminID,
		GameMaintenance:        false,
	}

	err := config.DB.FirstOrCreate(
		&maintenance,
		models.Maintenance{ID: 1},
	).Error

	if err != nil {
		return err
	}

	log.Println("Maintenance seeded successfully")

	return nil
}
