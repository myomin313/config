package Repositories

import (
	enums "config-service/Enums"
	"config-service/models"
	"config-service/requests"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BannerRepository struct {
	DB *gorm.DB
}

func NewBannerRepository(db *gorm.DB) *BannerRepository {
	return &BannerRepository{DB: db}
}

// Create inserts a new banner
func (r *BannerRepository) Create(banner *models.Banner, images []models.BannerImage) error {

	return r.DB.Transaction(func(tx *gorm.DB) error {

		// create banner
		if err := tx.Create(banner).Error; err != nil {
			return err
		}

		// attach banner_id
		for i := range images {
			images[i].BannerId = banner.ID
		}

		// insert images
		if len(images) > 0 {
			if err := tx.Create(&images).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// FindByID retrieves a banner by its ID
func (r *BannerRepository) FindByID(id uint64) (*models.Banner, error) {

	var banner models.Banner
	if err := r.DB.First(&banner, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("banner not found")
		}
		return nil, err
	}
	return &banner, nil
}

func (r *BannerRepository) FindBannerWithImagesById(id uint64) (*models.Banner, error) {

	var banner models.Banner
	if err := r.DB.Preload("BannerImages").First(&banner, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("banner not found")
		}
		return nil, err
	}
	return &banner, nil
}

// Update modifies an existing banner
func (r *BannerRepository) Update(banner *models.Banner, images []models.BannerImage) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {

		// Update banner
		if err := tx.Save(banner).Error; err != nil {
			return err
		}

		if len(images) > 0 {

			for i := range images {
				images[i].BannerId = banner.ID

				// 👉 Create new record
				if images[i].ID == 0 {
					if err := tx.Create(&images[i]).Error; err != nil {
						return err
					}
					continue
				}

				// 👉 Prepare update fields dynamically
				updateData := map[string]interface{}{
					"default_image": images[i].DefaultImage,
				}

				// Only update image_path if it's provided
				if images[i].ImagePath != nil {
					updateData["image_path"] = images[i].ImagePath
				}

				// 👉 Execute update
				if err := tx.Model(&models.BannerImage{}).
					Where("id = ? AND banner_id = ?", images[i].ID, banner.ID).
					Updates(updateData).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// Delete removes a banner by ID
func (r *BannerRepository) Delete(id uint64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Where("banner_id = ?", id).Delete(&models.BannerImage{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&models.Banner{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

// All fetches all banners
func (r *BannerRepository) All(offset int, limit int, req *requests.BannerSearchRequest) ([]models.Banner, int64, error) {
	var banners []models.Banner
	var total int64

	// Start building query
	query := r.DB.Model(&models.Banner{})
	// 🔍 Apply filters
	if req.ImageName != nil && *req.ImageName != "" {
		search := "%" + *req.ImageName + "%"
		query = query.Where("image_name LIKE ?", search)
	}
	if req.DisplayPlatform != nil {
		query = query.Where("display_platform = ?", enums.ParseDisplayPlatform(*req.DisplayPlatform))
	}

	if req.DisplayDuration != nil {
		query = query.Where("display_duration = ?", enums.ParseDisplayDuration(*req.DisplayDuration))
	}

	if req.Status != nil {
		query = query.Where("status = ?", enums.ParseStatus(*req.Status))
	}

	// ✅ Count total records for pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ✅ Apply pagination & order
	if err := query.Offset(offset).Limit(limit).Order("updated_at DESC").Find(&banners).Error; err != nil {
		return nil, 0, err
	}

	return banners, total, nil
}

func (r *BannerRepository) MassUpdateStatus(
	tableIDs []uint,
	status enums.Status,
) (int64, error) {

	result := r.DB.
		Model(&models.Banner{}).
		Where("id IN ?", tableIDs).
		Update("status", status)

	if result.Error != nil {
		return 0, result.Error
	}

	if result.RowsAffected == 0 {
		return 0, fmt.Errorf("no records updated")
	}

	return result.RowsAffected, nil
}

func (r *BannerRepository) MassDelete(
	tableIDs []uint,
) (int64, error) {

	result := r.DB.
		Model(&models.Banner{}).
		Where("id IN ?", tableIDs).
		Delete(&models.Banner{})

	if result.Error != nil {
		return 0, result.Error
	}

	if result.RowsAffected == 0 {
		return 0, fmt.Errorf("no_records_deleted")
	}

	return result.RowsAffected, nil
}

func (r *BannerRepository) GetActiveBanners(displayPlatform enums.DisplayPlatform, lang string) ([]models.Banner, error) {
	var banners []models.Banner
	now := time.Now()

	db := r.DB.
		Preload("BannerImages", func(tx *gorm.DB) *gorm.DB {
			return tx.Where("language_code = ? OR default_image = ?", lang, true)
		}).
		Where("display_platform = ?", displayPlatform).
		Where("status = ?", 1).
		Where(`
			(
				(display_duration = 1 AND start_date <= ? AND end_date >= ?)
				OR
				(display_duration = 2 AND start_date <= ?)
			)
		`, now, now, now)

	if err := db.Find(&banners).Error; err != nil {
		return nil, err
	}

	for i := range banners {
		var matched *models.BannerImage
		var fallback *models.BannerImage

		for j := range banners[i].BannerImages {
			img := &banners[i].BannerImages[j]

			if img.LanguageCode == lang {
				matched = img
				break
			}

			if img.DefaultImage != nil && *img.DefaultImage {
				fallback = img
			}
		}

		if matched != nil {
			banners[i].BannerImages = []models.BannerImage{*matched}
		} else if fallback != nil {
			banners[i].BannerImages = []models.BannerImage{*fallback}
		} else {
			banners[i].BannerImages = []models.BannerImage{}
		}
	}

	return banners, nil
}
func (r *BannerRepository) FindBannersWithImagesByIDs(ids []uint) ([]models.Banner, error) {
	var banners []models.Banner

	err := r.DB.
		Preload("BannerImages").
		Where("id IN ?", ids).
		Find(&banners).Error

	if err != nil {
		return nil, err
	}

	return banners, nil
}
