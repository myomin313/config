package Repositories

import (
	"config-service/models"

	"gorm.io/gorm"
)

type PlayTypeRepository struct {
	DB *gorm.DB
}

func (r *PlayTypeRepository) GetAll(play_data_id uint) ([]models.GameType, error) {
	var GameTypes []models.GameType

	err := r.DB.
		Where("game_data_id = ?", play_data_id).
		Order("id ASC").
		Find(&GameTypes).Error

	return GameTypes, err
}
func (r *PlayTypeRepository) FindByID(id uint64) (*models.GameType, error) {
	var playtype models.GameType

	if err := r.DB.First(&playtype, id).Error; err != nil {
		return nil, err
	}

	return &playtype, nil
}
