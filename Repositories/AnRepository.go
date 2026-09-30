package Repositories

import (
	"config-service/models"

	"log"

	"fmt"
	"time"

	"gorm.io/gorm"
)

type AnnouncementRepository struct {
	DB *gorm.DB
}

// Create a announcement
func (r *AnnouncementRepository) Create(announcement *models.Announcement) error {
	return r.DB.Create(announcement).Error
}

func (r *AnnouncementRepository) CreateAnnouncementDetail(announcementdetail *models.AnnouncementDetail) error {
	return r.DB.Create(announcementdetail).Error
}

// FindByID returns a announcement by ID
func (r *AnnouncementRepository) FindByID(id uint) (*models.Announcement, error) {
	var announcement models.Announcement

	if err := r.DB.
		Preload("AnnouncementDetail").
		First(&announcement, id).Error; err != nil {

		return nil, err
	}
	return &announcement, nil
}

// FindByID returns a announcement by ID
func (r *AnnouncementRepository) FindDetailByID(id uint) (*models.AnnouncementDetail, error) {
	var announcementdetail models.AnnouncementDetail

	if err := r.DB.
		First(&announcementdetail, id).Error; err != nil {

		return nil, err
	}
	return &announcementdetail, nil
}
func (r *AnnouncementRepository) LanguageExists(announcementID uint, languageCode string) (bool, error) {
	var count int64
	err := r.DB.Model(&models.AnnouncementDetail{}).
		Where("announcement_id = ? AND language_code = ?", announcementID, languageCode).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *AnnouncementRepository) LanguageExistsForUpdate(announcementID uint, languageCode string, ID uint) (bool, error) {
	var count int64
	err := r.DB.Model(&models.AnnouncementDetail{}).
		Where("announcement_id = ? AND language_code = ? AND id <> ?", announcementID, languageCode, ID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *AnnouncementRepository) Delete(announcement *models.Announcement) error {
	return r.DB.Delete(announcement).Error
}
func (r *AnnouncementRepository) DeleteDetail(announcementdetail *models.AnnouncementDetail) error {
	return r.DB.Delete(announcementdetail).Error
}
func (r *AnnouncementRepository) Update(announcement *models.Announcement) error {
	return r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(announcement).Error
}
func (r *AnnouncementRepository) UpdateDetail(announcementdetail *models.AnnouncementDetail) error {
	return r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(announcementdetail).Error
}

//	func (r *AnnouncementRepository) DeleteDetails(announcementID uint) error {
//		return r.DB.
//			Where("announcement_id = ?", announcementID).
//			Delete(&models.AnnouncementDetail{}).Error
//	}
func (r *AnnouncementRepository) FindAll(
	name, status, audience, language string,
	offset int,
	limit int,
) ([]models.Announcement, int64, error) {

	var announcements []models.Announcement
	var total int64

	query := r.DB.Model(&models.Announcement{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if audience != "" {
		query = query.Where("broadcast_audience = ?", audience)
	}

	if name != "" {
		query = query.Where("announcement_name LIKE ?", "%"+name+"%")
	}

	// ✅ Filter announcements by language (BUT NOT details)
	if language != "" {
		log.Println("Filter language:", language)

		subQuery := r.DB.Model(&models.AnnouncementDetail{}).
			Select("announcement_id").
			Where("language_code = ?", language)

		query = query.Where("id IN (?)", subQuery)
	}

	// ✅ Count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ✅ Always preload ALL details (important!)
	err := query.
		Preload("AnnouncementDetail").
		Offset(offset).
		Limit(limit).
		Order("sort_order ASC").
		Find(&announcements).Error

	return announcements, total, err
}
func (r *AnnouncementRepository) FindAllForFE(broadcastAudience, lang string) ([]models.Announcement, error) {
	var announcements []models.Announcement

	query := r.DB.Model(&models.Announcement{}).
		Where("status = ?", "enable").
		Where("broadcast_audience IN ?", []string{broadcastAudience, "all"}). // ✅ matches audience or "all"
		Preload("AnnouncementDetail", "language_code = ?", lang).
		Order("sort_order ASC")

	err := query.Find(&announcements).Error
	return announcements, err
}
func (r *AnnouncementRepository) MassUpdateStatus(AnnouncementIDs []uint, status string) error {

	result := r.DB.
		Model(&models.Announcement{}).
		Where("id IN ?", AnnouncementIDs).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no_records_updated")
	}
	return nil
}
func (r *AnnouncementRepository) DeleteMultiple(ids []uint) error {
	return r.DB.
		Where("id IN ?", ids).
		Delete(&models.Announcement{}).Error
}

func (r *AnnouncementRepository) UpdateAnnouncementTimestamp(id uint) error {
	return r.DB.Model(&models.Announcement{}).
		Where("id = ?", id).
		Update("updated_at", time.Now()).Error
}
func (r *AnnouncementRepository) FindByIDs(ids []uint) ([]models.Announcement, error) {
	var announcements []models.Announcement

	err := r.DB.
		Preload("AnnouncementDetail").
		Where("id IN ?", ids).
		Find(&announcements).Error

	if err != nil {
		return nil, err
	}

	return announcements, nil
}
