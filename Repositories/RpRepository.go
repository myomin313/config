package Repositories

import (
	"config-service/models"

	"gorm.io/gorm"
)

type ReportRepository struct {
	DB *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{DB: db}
}

// All fetches all banners

func (r *ReportRepository) GetPlayResult(offset int, limit int) ([]models.Game, int64, error) {
	var playtables []models.Game
	var total int64

	// Count total records
	if err := r.DB.Model(&models.Game{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	//Get paginated data
	err := r.DB.
		//Preload("GameTable").
		Preload("GameTypes").
		Preload("Translations", "lang_name = ?", "game_name").
		Preload("GameTypeData.MasterGameType.GameTypeTranslation").
		Offset(offset).
		Limit(limit).
		Find(&playtables).Error

	if err != nil {
		return nil, 0, err
	}

	return playtables, total, nil
}

func (r *ReportRepository) GetPlayNameByPlayID(playID uint64, languageCode string) (string, error) {
	var play models.Game

	err := r.DB.
		Preload("GameTypes").
		Preload("Translations", "language_code = ?", languageCode).
		Preload("GameTypeData.MasterGameType.GameTypeTranslation").
		Where("id = ?", playID).
		First(&play).Error

	if err != nil {
		return "", err
	}

	if len(play.Translations) > 0 {
		return play.Translations[0].Content, nil
	}

	return "", nil
}

func (r *ReportRepository) CreateRoundRemarkLog(log *models.RoundRemarkLog) error {
	return r.DB.Create(log).Error
}

func (r *ReportRepository) CreateRoundResultModificationLog(log *models.RoundResultModificationLog) error {
	return r.DB.Create(log).Error
}
