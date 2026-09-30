package Repositories

import (
	"config-service/models"

	"gorm.io/gorm"
)

type PlatformApiKeyRepository struct {
	DB *gorm.DB
}

func (r *PlatformApiKeyRepository) Create(apiKey *models.PlatformApiKey) error {
	return r.DB.Create(apiKey).Error
}

func (r *PlatformApiKeyRepository) FindByPlatformID(
	platformID uint64,
	apiKey *models.PlatformApiKey,
) error {
	return r.DB.Where("platform_id = ?", platformID).First(apiKey).Error
}

func (r *PlatformApiKeyRepository) Update(apiKey *models.PlatformApiKey) error {
	return r.DB.Save(apiKey).Error
}

func (r *PlatformApiKeyRepository) Delete(apiKey *models.PlatformApiKey) error {
	return r.DB.Delete(apiKey).Error
}
