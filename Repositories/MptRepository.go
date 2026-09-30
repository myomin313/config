package Repositories

import (
	"config-service/models"

	"gorm.io/gorm"
)

type MasterPlayTypeRepository struct {
	DB *gorm.DB
}

// func (r *MasterGameTypeRepository) GetAll(gameTypeID uint64)   ([]models.MasterGameType, error) {
// 	var masterGameTypes []models.MasterGameType

// 	err := r.DB.
// 	    Where("game_type_id = ?", gameTypeID).
// 	    Preload("GameTypeTranslation").
// 		Preload("GameType").
// 		Order("id ASC").
// 		Find(&masterGameTypes).Error

// 	return masterGameTypes, err
// }

func (r *MasterPlayTypeRepository) GetAll(playTypeIDs ...uint64) ([]models.MasterGameType, error) {
	var masterPlayTypes []models.MasterGameType

	query := r.DB.Preload("GameTypeTranslation").Preload("GameType").Order("id ASC")

	if len(playTypeIDs) > 0 {
		query = query.Where("game_type_id IN ?", playTypeIDs)
	}

	err := query.Find(&masterPlayTypes).Error
	return masterPlayTypes, err
}
