package Repositories

import (
	"config-service/config"
	"config-service/models"
	"fmt"

	"gorm.io/gorm"
)

type PlayHallRepository struct {
	DB *gorm.DB
}

func (r *PlayHallRepository) Create(playhall *models.GameHall) error {
	return r.DB.Create(playhall).Error
}

func (r *PlayHallRepository) FindByID(id uint64) (*models.GameHall, error) {
	var playhall models.GameHall

	if err := r.DB.
		Preload("Translations", "lang_name = ?", "game_hall_name").
		First(&playhall, id).Error; err != nil {

		return nil, err
	}

	return &playhall, nil
}
func (r *PlayHallRepository) FindAll() ([]models.GameHall, error) {
	var playhalls []models.GameHall

	err := r.DB.Order("sort_order ASC").
		Preload("Translations", "lang_name = ?", "game_hall_name"). // preload master for names
		Find(&playhalls).Error
	if err != nil {
		config.Error(fmt.Sprintf("Error fetching game halls: %v", err))
	}

	return playhalls, err
}
func (r *PlayHallRepository) FindAllEnable() ([]models.GameHall, error) {
	var playhalls []models.GameHall

	err := r.DB.
		Where("status = ? ", "enable").
		Preload("Translations", "lang_name = ?", "game_hall_name"). // preload master for names
		Find(&playhalls).Error

	return playhalls, err
}
func (r *PlayHallRepository) Delete(playhall *models.GameHall) error {
	return r.DB.Delete(playhall).Error
}

// Update game with associations (translations + GameTypeData)
func (r *PlayHallRepository) Update(playhall *models.GameHall) error {
	return r.DB.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(playhall).Error
}
