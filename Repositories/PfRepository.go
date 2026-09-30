package Repositories

import (
	"config-service/models"

	"gorm.io/gorm"
)

// PlatformRepository handles DB operations for Platform
type PlatformRepository struct {
	DB *gorm.DB
}

// FindAll returns all Platforms
func (r *PlatformRepository) FindAll(platforms *[]models.Platform) error {
	//return r.DB.Preload("ApiKey").Find(platforms).Error
	return r.DB.Preload("APIKey").Find(platforms).Error
}
func (r *PlatformRepository) FindAllEnable(platforms *[]models.Platform) error {
	//return r.DB.Preload("ApiKey").Find(platforms).Error
	return r.DB.Where("status = ?", "enable").Find(platforms).Error
}

// FindByID returns a Platform by ID
func (r *PlatformRepository) FindByID(id uint64, platform *models.Platform) error {
	return r.DB.Preload("APIKey").First(platform, id).Error
}

// Update updates an existing Platform
func (r *PlatformRepository) Update(platform *models.Platform) error {
	return r.DB.Save(platform).Error
}

// Delete deletes a Platform
func (r *PlatformRepository) Delete(platform *models.Platform) error {
	return r.DB.Delete(platform).Error
}
func (r *PlatformRepository) GetMerchantCurrency(currencies *[]models.MerchantCurrency) error {
	return r.DB.Find(currencies).Error
}
func (r *PlatformRepository) Exists(
	merchantCode string,
) (bool, error) {

	var count int64

	err := r.DB.
		Model(&models.MerchantCurrency{}).
		Where("merchant_code = ?", merchantCode).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
func (r *PlatformRepository) CreateMerchantCurrency(
	merchantCurrency *models.MerchantCurrency,
) error {
	return r.DB.Create(merchantCurrency).Error
}

func (r *PlatformRepository) CreatePlatformWithApiKey(
	platform *models.Platform,
	apiKey *models.PlatformApiKey,
) error {

	tx := r.DB.Begin()

	if err := tx.Create(platform).Error; err != nil {
		tx.Rollback()
		return err
	}

	apiKey.PlatformID = platform.ID

	if err := tx.Create(apiKey).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
func (r *PlatformRepository) GetApiKeyByPlatformID(platformID uint64) (*models.PlatformApiKey, error) {
	var apiKey models.PlatformApiKey

	if err := r.DB.
		Where("platform_id = ?", platformID).
		First(&apiKey).Error; err != nil {
		return nil, err
	}

	return &apiKey, nil
}

func (r *PlatformRepository) UpdateApiKey(apiKey *models.PlatformApiKey) error {
	return r.DB.Save(apiKey).Error
}
func (r *PlatformRepository) GetMerchantCurrencyByPlatformID(
	platformID uint64,
	currencies *[]models.MerchantCurrency,
) error {
	return r.DB.
		Preload("Platform").
		Where("platform_id = ?", platformID).
		Find(currencies).Error
}
func (r *PlatformRepository) FindMerchantCurrencyByID(id uint) (*models.MerchantCurrency, error) {
	var merchantCurrency models.MerchantCurrency

	if err := r.DB.First(&merchantCurrency, id).Error; err != nil {
		return nil, err
	}

	return &merchantCurrency, nil
}

func (r *PlatformRepository) UpdateMerchantCurrency(
	merchantCurrency *models.MerchantCurrency,
) error {
	return r.DB.Save(merchantCurrency).Error
}
func (r *PlatformRepository) ExistsMerchantCurrency(
	merchantCode string,
	excludeID uint,
) (bool, error) {

	var count int64

	err := r.DB.Model(&models.MerchantCurrency{}).
		Where("merchant_code = ? AND id <> ?", merchantCode, excludeID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
