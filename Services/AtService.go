package Services

import (
	"config-service/Repositories"
	"config-service/helper"
	"config-service/models"
	"config-service/requests"
	"time"

	// "gorm.io/gorm"
	"errors"
	"fmt"

	//"log"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AnnouncementService struct {
	Repo  *Repositories.AnnouncementRepository
	Audit *helper.AuditService
}

func (s *AnnouncementService) CreateAnnouncement(req requests.AnnouncementCreateRequest, id uint) (*models.Announcement, error) {

	announcement := models.Announcement{
		AnnouncementName:  req.AnnouncementName,
		SortOrder:         req.SortOrder,
		Status:            req.Status,
		BroadcastAudience: req.BroadcastAudience,
		CreatedBy:         id,
		UpdatedBy:         id,
	}
	// for _, gtd := range req.AnnouncementDetail { // make sure your request struct has GameTypeDataRequest slice
	// 	AnnouncementDetail := models.AnnouncementDetail{
	// 		LanguageCode:     gtd.LanguageCode,
	// 		Title:            gtd.Title,
	// 		Content:          gtd.Content,
	// 		CreatedAt:        time.Now(),
	// 		UpdatedAt:        time.Now(),
	// 	}
	// 	announcement.AnnouncementDetail = append(announcement.AnnouncementDetail, AnnouncementDetail)
	// }
	if err := s.Repo.Create(&announcement); err != nil {
		return nil, err
	}
	return s.Repo.FindByID(announcement.ID)
}
func (s *AnnouncementService) CreateAnnouncementDetail(req requests.AnnouncementDetailRequest, AnnouncementID uint) (*models.AnnouncementDetail, error) {
	announcement, err := s.Repo.FindByID(AnnouncementID)
	if err != nil {
		return nil, err
	}
	exists, err := s.Repo.LanguageExists(AnnouncementID, req.LanguageCode)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("language already exists for this announcement")
	}
	announcementdetail := models.AnnouncementDetail{
		LanguageCode:   req.LanguageCode,
		Title:          req.Title,
		Content:        req.Content,
		AnnouncementID: announcement.ID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.Repo.CreateAnnouncementDetail(&announcementdetail); err != nil {
		return nil, err
	}
	// Update parent announcement UpdatedAt
	if err := s.Repo.UpdateAnnouncementTimestamp(AnnouncementID); err != nil {
		return nil, err
	}
	return s.Repo.FindDetailByID(announcementdetail.ID)
}
func (s *AnnouncementService) Delete(id uint) error {
	announcement, err := s.Repo.FindByID(id)
	if err != nil {
		return err
	}
	if err := s.Repo.Delete(announcement); err != nil {
		return fmt.Errorf("failed_to_delete_announcement: %w", err)
	}
	return nil
}
func (s *AnnouncementService) DeleteDetail(id uint) error {
	announcementdetail, err := s.Repo.FindDetailByID(id)
	if err != nil {
		return err
	}
	if err := s.Repo.DeleteDetail(announcementdetail); err != nil {
		return fmt.Errorf("failed_to_delete_announcement_detail: %w", err)
	}
	if err := s.Repo.UpdateAnnouncementTimestamp(announcementdetail.AnnouncementID); err != nil {
		return err
	}
	return nil
}
func (s *AnnouncementService) Update(AnnouncementID uint, req requests.AnnouncementUpdateRequest, id uint) (*models.Announcement, error) {
	announcement, err := s.Repo.FindByID(AnnouncementID)
	if err != nil {
		return nil, err
	}
	announcement.AnnouncementName = req.AnnouncementName
	announcement.SortOrder = req.SortOrder
	announcement.Status = req.Status
	announcement.BroadcastAudience = req.BroadcastAudience
	announcement.UpdatedBy = id
	announcement.UpdatedAt = time.Now()
	// Delete old details
	// err = s.Repo.DeleteDetails(AnnouncementID)
	// if err != nil {
	// 	return nil, err
	// }
	// announcement.AnnouncementDetail = []models.AnnouncementDetail{}
	// // Game types
	// for _, gtd := range req.AnnouncementDetail { // make sure your request struct has GameTypeDataRequest slice
	// 	AnnouncementDetail := models.AnnouncementDetail{
	// 		AnnouncementID:   AnnouncementID,
	// 		LanguageCode:     gtd.LanguageCode,
	// 		Title:            gtd.Title,
	// 		Content:          gtd.Content,
	// 		CreatedAt:        time.Now(),
	// 		UpdatedAt:        time.Now(),
	// 	}
	// 	// Link to parent game after we have ID
	// 	// For GORM, we can use association later, but here we append to AnnouncementDetail slice
	// 	//announcement.AnnouncementDetail = append(AnnouncementDetail)
	// 	announcement.AnnouncementDetail = append(announcement.AnnouncementDetail, AnnouncementDetail)
	// }
	if err := s.Repo.Update(announcement); err != nil {
		return nil, err
	}
	return s.Repo.FindByID(announcement.ID)
}
func (s *AnnouncementService) UpdateAnnouncementDetail(ID uint, req requests.AnnouncementDetailRequest, id uint) (*models.AnnouncementDetail, error) {
	announcementdetail, err := s.Repo.FindDetailByID(ID)
	if err != nil {
		return nil, err
	}
	//Check uniqueness of language
	exists, err := s.Repo.LanguageExistsForUpdate(announcementdetail.AnnouncementID, req.LanguageCode, ID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("language already exists for this announcement")
	}
	announcementdetail.LanguageCode = req.LanguageCode
	announcementdetail.Title = req.Title
	announcementdetail.Content = req.Content
	announcementdetail.UpdatedAt = time.Now()

	if err := s.Repo.UpdateDetail(announcementdetail); err != nil {
		return nil, err
	}
	if err := s.Repo.UpdateAnnouncementTimestamp(announcementdetail.AnnouncementID); err != nil {
		return nil, err
	}
	return s.Repo.FindDetailByID(ID)
}
func (s *AnnouncementService) GetAll(
	c *fiber.Ctx,
	req *requests.AnnouncementSearchRequest,
) ([]models.Announcement, int64, int, int, error) {

	pagination := helper.GetPagination(c)

	var name, status, audience, language string

	if req.AnnouncementName != nil {
		name = *req.AnnouncementName
	}
	if req.Status != nil {
		status = *req.Status
	}
	if req.BroadcastAudience != nil {
		audience = *req.BroadcastAudience
	}
	if req.LanguageCode != nil {
		language = *req.LanguageCode
	}

	announcements, total, err := s.Repo.FindAll(
		name,
		status,
		audience,
		language,
		pagination.Offset,
		pagination.Limit,
	)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	return announcements, total, pagination.PageNo, pagination.PageSize, nil
}
func (s *AnnouncementService) GetAllForFE(broadcastAudience, lang string) ([]models.Announcement, error) {
	return s.Repo.FindAllForFE(broadcastAudience, lang)
}

// update the Announcement status
func (s *AnnouncementService) MassUpdateStatus(AnnouncementIDs []uint, status string) error {
	return s.Repo.MassUpdateStatus(AnnouncementIDs, status)
}
func (s *AnnouncementService) DeleteMultiple(ids []uint) error {
	return s.Repo.DeleteMultiple(ids)
}
func (s *AnnouncementService) GetAnnouncementsDetail(id uint) (*models.Announcement, error) {
	return s.Repo.FindByID(id)
}
func (s *AnnouncementService) AuditAction(
	c *fiber.Ctx,
	action string,
	model string,
	user jwt.MapClaims,
	oldData any,
	newData any,
) {

	helper.PublishAudit(
		s.Audit,
		c,
		action,
		model,
		user,
		oldData,
		newData,
	)
}
func (s *AnnouncementService) GetAnnouncementByID(id uint) (*models.Announcement, error) {
	return s.Repo.FindByID(id)
}
func (s *AnnouncementService) GetAnnouncementDetailByID(id uint) (*models.AnnouncementDetail, error) {
	return s.Repo.FindDetailByID(id)
}
func (s *AnnouncementService) GetAnnouncementsByIDs(ids []uint) ([]models.Announcement, error) {
	return s.Repo.FindByIDs(ids)
}
