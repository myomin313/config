package Repositories

import (
	"config-service/models"

	"gorm.io/gorm"
)

type PlayDataRepository struct {
	DB *gorm.DB
}

func (r *PlayDataRepository) FindAll() ([]models.GameData, error) {
	var playdata []models.GameData

	err := r.DB.Find(&playdata).Error

	return playdata, err
}
