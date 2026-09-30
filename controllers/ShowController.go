package controllers

import (
	//"context"
	"config-service/ApiResponse"
	"config-service/Repositories"
	Resources "config-service/Resources"
	"config-service/Services"
	"config-service/config"
	"config-service/helper"
	"config-service/models"
	"config-service/requests"
	"config-service/utils"
	"log"
	"math"

	//	"os"
	//"sort"
	"time"

	"config-service/pb/commonpb"

	"context"

	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	//"fmt"
	//"errors"
)

// ------------------------ Controller Struct ------------------------
type AnnouncementController struct {
	Service *Services.AnnouncementService
}

func NewAnnouncementController(db *gorm.DB) *AnnouncementController {
	auditService := &helper.AuditService{
		//Redis: redisConn,
		Channel: config.RabbitChannel,
		Queue:   "audit_logs",
	}
	return &AnnouncementController{
		Service: &Services.AnnouncementService{
			Repo: &Repositories.AnnouncementRepository{
				DB: db,
			},
			Audit: auditService,
		},
	}
}

func mapAnnouncementToResouce(announcement *models.Announcement) Resources.AnnouncementResource {

	conn := config.GrpcClient()
	defer conn.Close()

	client := commonpb.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Call GetUser
	resp, err := client.GetUser(ctx, &commonpb.GetUserRequest{Id: int64(announcement.CreatedBy)})
	if err != nil {
		log.Println("Error calling GetUser:", err)
	}
	respupdate, err := client.GetUser(ctx, &commonpb.GetUserRequest{Id: int64(announcement.UpdatedBy)})
	if err != nil {
		log.Println("Error calling GetUser:", err)
	}
	var details []Resources.AnnouncementDetailResource

	for _, d := range announcement.AnnouncementDetail {
		details = append(details, Resources.AnnouncementDetailResource{
			ID:           d.ID,
			LanguageCode: d.LanguageCode,
			Title:        d.Title,
			Content:      d.Content,
		})
	}

	// title := ""

	// for _, d := range announcement.AnnouncementDetail {
	// 	if d.LanguageCode == lang {
	// 		title = d.Title
	// 		break
	// 	}
	// }

	return Resources.AnnouncementResource{
		ID:               announcement.ID,
		AnnouncementName: announcement.AnnouncementName,
		//Title:       title,
		Status:             announcement.Status,
		SortOrder:          announcement.SortOrder,
		BroadcastAudience:  announcement.BroadcastAudience,
		CreatedBy:          resp.Name,
		UpdatedBy:          respupdate.Name,
		CreatedAt:          announcement.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:          announcement.UpdatedAt.Format("2006-01-02 15:04:05"),
		AnnouncementDetail: details,
	}

}

func mapAnnouncementFEToResouce(announcement *models.Announcement, lang string) Resources.AnnouncementListFEResource {
	title := ""
	content := ""
	for _, d := range announcement.AnnouncementDetail {
		if d.LanguageCode == lang {
			title = d.Title
			content = d.Content
			break
		}
	}

	return Resources.AnnouncementListFEResource{
		ID:                announcement.ID,
		Title:             title,
		Content:           content,
		Status:            announcement.Status,
		SortOrder:         announcement.SortOrder,
		BroadcastAudience: announcement.BroadcastAudience,
		CreatedAt:         announcement.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         announcement.UpdatedAt.Format("2006-01-02 15:04:05"),
		// If you need title, add it to FE struct
		// Title: title,
	}
}
func mapAnnouncementLogToResouce(announcement *models.Announcement) Resources.AnnouncementLogResource {

	conn := config.GrpcClient()
	defer conn.Close()

	client := commonpb.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Call GetUser
	resp, err := client.GetUser(ctx, &commonpb.GetUserRequest{Id: int64(announcement.CreatedBy)})
	if err != nil {
		log.Println("Error calling GetUser:", err)
	}
	respupdate, err := client.GetUser(ctx, &commonpb.GetUserRequest{Id: int64(announcement.UpdatedBy)})
	if err != nil {
		log.Println("Error calling GetUser:", err)
	}
	return Resources.AnnouncementLogResource{
		ID:               announcement.ID,
		AnnouncementName: announcement.AnnouncementName,
		//Title:       title,
		Status:            announcement.Status,
		SortOrder:         announcement.SortOrder,
		BroadcastAudience: announcement.BroadcastAudience,
		CreatedBy:         resp.Name,
		UpdatedBy:         respupdate.Name,
		CreatedAt:         announcement.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         announcement.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

}
func mapAnnouncementDetailLogToResouce(announcementDetail *models.AnnouncementDetail, announcement *models.Announcement) Resources.AnnouncementDetailLogResource {

	return Resources.AnnouncementDetailLogResource{
		ID:               announcementDetail.ID,
		AnnouncementName: announcement.AnnouncementName,
		Title:            announcementDetail.Title,
		LanguageCode:     announcementDetail.LanguageCode,
		Content:          announcementDetail.Content,
		CreatedAt:        announcementDetail.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        announcementDetail.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

}

// func (gc *AnnouncementController) CreateAnnouncement(c *fiber.Ctx) error {
// 	req := new(requests.AnnouncementCreateRequest)
//     //lang := c.Get("Accept-Language", "en")
// 	user, _ := c.Locals("user").(jwt.MapClaims)

// 	if err := c.BodyParser(req); err != nil {
// 		fullMessages := utils.ValidationErrorMessages(err, req)
// 		return ApiResponse.ValidationError(c, fiber.StatusBadRequest, "unable_to_parse_request_body", fullMessages[0])
// 	}

// 	if err := utils.Validate.Struct(req); err != nil {
// 		fullMessages := utils.ValidationErrorMessages(err, req)
// 		log.Println("err", fullMessages[0])
// 		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation_error", fullMessages[0])
// 	}
// 	//game, err := gc.Service.CreateGame(*req)
// 	announcement, err := gc.Service.CreateAnnouncement(*req, uint(user["id"].(float64)))
// 	if err != nil {
// 		return ApiResponse.Error(c, fiber.StatusInternalServerError,err.Error())
// 	}
// 	resp := mapAnnouncementToResouce(announcement)
// 	return ApiResponse.Success(c, fiber.StatusCreated, resp, "announcement_created_successfully")
// }

// CreateAnnouncement godoc
// @Summary      Create a Announcement
// @Description  Create a new announcement
// @Tags         Announcements
// @Accept       json
// @Produce      json
// @Param        payload body requests.AnnouncementCreateRequest true "Announcement Data"
// @Success      201 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /announcements [post]
func (gc *AnnouncementController) CreateAnnouncement(c *fiber.Ctx) error {
	req := new(requests.AnnouncementCreateRequest)
	//lang := c.Get("Accept-Language", "en")
	user, _ := c.Locals("user").(jwt.MapClaims)

	if err := c.BodyParser(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		return ApiResponse.ValidationError(c, fiber.StatusBadRequest, "unable_to_parse_request_body", fullMessages[0])
	}

	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("err", fullMessages[0])
		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation_error", fullMessages[0])
	}
	//game, err := gc.Service.CreateGame(*req)
	announcement, err := gc.Service.CreateAnnouncement(*req, uint(user["id"].(float64)))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	announcementlog := mapAnnouncementLogToResouce(announcement)
	gc.Service.AuditAction(
		c,
		"create",
		"announcement",
		user,
		nil,
		announcementlog,
	)
	//resp := mapAnnouncementToResouce(announcement)
	return ApiResponse.Success(c, fiber.StatusCreated, announcement, "announcement_created_successfully")
}

// CreateAnnouncementDetail godoc
// @Summary      Create a Announcement Detail
// @Description  Create a new announcement
// @Tags         Announcements
// @Accept       json
// @Produce      json
// @Param        payload body requests.AnnouncementDetailRequest true "Announcement Detail Data"
// @Success      201 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /announcements/{id}/detail [post]
func (gc *AnnouncementController) CreateAnnouncementDetail(c *fiber.Ctx) error {
	req := new(requests.AnnouncementDetailRequest)
	//lang := c.Get("Accept-Language", "en")
	//user, _ := c.Locals("user").(jwt.MapClaims)
	AnnouncementID, err := c.ParamsInt("id")

	if err := c.BodyParser(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		return ApiResponse.ValidationError(c, fiber.StatusBadRequest, "unable_to_parse_request_body", fullMessages[0])
	}

	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("err", fullMessages[0])
		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation_error", fullMessages[0])
	}
	//game, err := gc.Service.CreateGame(*req)
	announcement, err := gc.Service.CreateAnnouncementDetail(*req, uint(AnnouncementID))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	user, _ := c.Locals("user").(jwt.MapClaims)
	announcementDetail, err := gc.Service.GetAnnouncementByID(uint(announcement.AnnouncementID))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	log := mapAnnouncementDetailLogToResouce(announcement, announcementDetail)
	gc.Service.AuditAction(
		c,
		"create",
		"announcement detail",
		user,
		nil,
		log,
	)
	//resp := mapAnnouncementToResouce(announcement)
	return ApiResponse.Success(c, fiber.StatusCreated, announcement, "announcement_detail_created_successfully")
}

// DeleteAnnouncement godoc
// @Summary      Delete announcements
// @Description  Delete announcements by ID
// @Tags         Announcements
// @Param        id path int true "announcementsID"
// @Success      204 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /announcements/{id} [delete]
func (gc *AnnouncementController) DeleteAnnouncement(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	log.Println("hello")
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid_id")
	}
	if id < 0 {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid_id")
	}
	deleteAnnouncement, _ := gc.Service.GetAnnouncementByID(uint(id))
	err = gc.Service.Delete(uint(id))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	user, _ := c.Locals("user").(jwt.MapClaims)

	announcementlog := mapAnnouncementLogToResouce(deleteAnnouncement)
	gc.Service.AuditAction(
		c,
		"delete",
		"announcement",
		user,
		announcementlog,
		nil,
	)

	// Return empty success response for consistency
	return ApiResponse.Success(c, 200, map[string]interface{}{}, "announcement_deleted_successfully")
}

// UpdateAnnouncement godoc
// @Summary      Update announcements
// @Description  Update announcements by ID
// @Tags         Announcements
// @Accept       json
// @Produce      json
// @Param        id path int true "Announcement ID"
// @Param        payload body requests.AnnouncementUpdateRequest true "Updated announcement Data"
// @Success      200 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /announcements/{id} [put]
func (gc *AnnouncementController) UpdateAnnouncement(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	//lang := c.Get("Accept-Language", "en")
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid_id")
	}
	req := new(requests.AnnouncementUpdateRequest)
	if err := c.BodyParser(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("bad request", err.Error())
		return ApiResponse.ValidationError(c, fiber.StatusBadRequest, "unable_to_parse_request_body", fullMessages[0])
	}

	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("update err", fullMessages[0])
		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation_error", fullMessages[0])
	}

	user, _ := c.Locals("user").(jwt.MapClaims)
	if id < 0 {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid_id")
	}
	userID := user["id"].(float64)
	if userID < 0 || userID > float64(math.MaxUint) {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid_user_id")
	}
	oldAnnouncement, _ := gc.Service.GetAnnouncementByID(uint(id))

	announcements, errNew := gc.Service.Update(uint(id), *req, uint(userID))

	if errNew != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, errNew.Error())
	}

	resp := mapAnnouncementToResouce(announcements)

	oldresp := mapAnnouncementToResouce(oldAnnouncement)

	// exclude AnnouncementDetail from oldresp
	oldData := map[string]interface{}{
		"announcement_name":  oldresp.AnnouncementName,
		"status":             oldresp.Status,
		"sort_order":         oldresp.SortOrder,
		"broadcast_audience": oldresp.BroadcastAudience,
		"created_by":         oldresp.CreatedBy,
		"updated_by":         oldresp.UpdatedBy,
		"created_at":         oldresp.CreatedAt,
		"updated_at":         oldresp.UpdatedAt,
	}

	newData := map[string]interface{}{
		"announcement_name":  resp.AnnouncementName,
		"status":             resp.Status,
		"sort_order":         resp.SortOrder,
		"broadcast_audience": resp.BroadcastAudience,
		"created_by":         resp.CreatedBy,
		"updated_by":         resp.UpdatedBy,
		"created_at":         resp.CreatedAt,
		"updated_at":         resp.UpdatedAt,
	}

	gc.Service.AuditAction(
		c,
		"update",
		"announcement",
		user,
		oldData,
		newData,
	)

	return ApiResponse.Success(c, 200, resp, "announcement_updated_successfully")
}

// UpdateAnnouncementDetail godoc
// @Summary      Update announcements detail
// @Description  Update announcements by ID
// @Tags         Announcements
// @Accept       json
// @Produce      json
// @Param        id path int true "Announcement ID"
// @Param        payload body requests.AnnouncementDetailRequest true "Updated announcement Data"
// @Success      200 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /announcements/detail/update/{id} [put]
func (gc *AnnouncementController) UpdateAnnouncementDetail(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	//lang := c.Get("Accept-Language", "en")
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid_id")
	}
	req := new(requests.AnnouncementDetailRequest)
	if err := c.BodyParser(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("bad request", err.Error())
		return ApiResponse.ValidationError(c, fiber.StatusBadRequest, "unable_to_parse_request_body", fullMessages[0])
	}

	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("update err", fullMessages[0])
		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation_error", fullMessages[0])
	}

	user, _ := c.Locals("user").(jwt.MapClaims)
	if id < 0 {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid_id")
	}
	userID := user["id"].(float64)
	if userID < 0 || userID > float64(math.MaxUint) {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid_user_id")
	}
	oldAnnouncementDetail, _ := gc.Service.GetAnnouncementDetailByID(uint(id))
	log.Println("oldAnnouncementDetail", oldAnnouncementDetail)

	announcementold, err := gc.Service.GetAnnouncementByID(uint(oldAnnouncementDetail.AnnouncementID))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	oldlog := mapAnnouncementDetailLogToResouce(oldAnnouncementDetail, announcementold)

	announcementdetails, errNew := gc.Service.UpdateAnnouncementDetail(uint(id), *req, uint(userID))
	log.Println("announcementdetails", announcementdetails)
	if errNew != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, errNew.Error())
	}
	announcementnew, err := gc.Service.GetAnnouncementByID(uint(announcementdetails.AnnouncementID))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	newlog := mapAnnouncementDetailLogToResouce(announcementdetails, announcementnew)

	gc.Service.AuditAction(
		c,
		"update",
		"announcement detail",
		user,
		oldlog,
		newlog,
	)

	//resp := mapAnnouncementToResouce(announcements)
	return ApiResponse.Success(c, 200, announcementdetails, "announcement_detail_updated_successfully")
}

// GetAnnouncementsList godoc
// @Summary      List Announcements
// @Description  Get announcements with pagination and filters
// @Tags         Announcements
// @Produce      json
// @Param        title                query     string  false  "Search by title"
// @Param        status               query     string  false  "Filter by status"
// @Param        broadcast_audience   query     string  false  "Filter by audience"
// @Param        page                 query     int     false  "Page number" default(1)
// @Param        page_size            query     int     false  "Page size" default(10)
// @Param        Accept-Language      header    string  false  "Language code (en, km)" default(en)
// @Success      200 {object} ApiResponse.Response "Successful response with pagination"
// @Failure      400 {object} ApiResponse.Response "Invalid query params"
// @Failure      500 {object} ApiResponse.Response "Internal server error"
// @Router       /announcements [get]
func (gc *AnnouncementController) GetAnnouncementsList(c *fiber.Ctx) error {
	req := new(requests.AnnouncementSearchRequest)
	// ✅ parse query into struct
	if err := c.QueryParser(req); err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid query params")
	}
	//lang := c.Get("Accept-Language", "en")

	announcements, total, pageNo, pageSize, err :=
		gc.Service.GetAll(c, req)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	// ✅ mapping stays here
	resp := make([]Resources.AnnouncementResource, 0, len(announcements))
	for i := range announcements {
		resp = append(resp, mapAnnouncementToResouce(&announcements[i]))
	}
	totalPage := (total + int64(pageSize) - 1) / int64(pageSize)

	return ApiResponse.Pagination(c, fiber.StatusOK, map[string]interface{}{
		"records":    resp,
		"page_no":    pageNo,
		"page_size":  pageSize,
		"total":      total,
		"total_page": totalPage,
	}, "announcements_retrieved_successfully")
}

// MassUpdateAnnouncementStatus godoc
// @Summary      Update  announcements status
// @Description  Mass update announcements status
// @Tags         Announcements
// @Param        payload body requests.GameTableMassStatusUpdateRequest true "Update announcements status"
// @Success      200 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /announcements/status/{ids} [post]
func (pc *AnnouncementController) MassUpdateAnnouncementStatus(c *fiber.Ctx) error {
	req := new(requests.AnnouncementMassStatusUpdateRequest)
	if err := c.BodyParser(&req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("bad request", err.Error())
		return ApiResponse.ValidationError(c, fiber.StatusBadRequest, "unable_to_parse_request_body", fullMessages[0])
	}

	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("update err", fullMessages[0])
		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation_error", fullMessages[0])
	}
	oldAnnouncements, err := pc.Service.GetAnnouncementsByIDs(req.AnnouncementIDs)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	err = pc.Service.MassUpdateStatus(
		req.AnnouncementIDs,
		req.Status,
	)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	newAnnouncements, err := pc.Service.GetAnnouncementsByIDs(req.AnnouncementIDs)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	oldItems := make([]map[string]interface{}, 0)

	for _, a := range oldAnnouncements {
		oldItems = append(oldItems, map[string]interface{}{
			"announcement_name": a.AnnouncementName,
			"status":            a.Status,
		})
	}

	newItems := make([]map[string]interface{}, 0)

	for _, a := range newAnnouncements {
		newItems = append(newItems, map[string]interface{}{
			"announcement_name": a.AnnouncementName,
			"status":            a.Status,
		})
	}

	oldData := map[string]interface{}{
		"announcement": oldItems,
	}

	newData := map[string]interface{}{
		"announcement": newItems,
	}
	pc.Service.AuditAction(
		c,
		"mass update",
		"announcement status",
		c.Locals("user").(jwt.MapClaims),
		oldData,
		newData,
	)

	return ApiResponse.Success(c, 200, len(req.AnnouncementIDs), "announcements_updated_successfully")
}

// func (gc *AnnouncementController) DeleteMultipleAnnouncement(c *fiber.Ctx) error {
// 	var req requests.DeleteAnnouncementsRequest

// 	if err := c.BodyParser(&req); err != nil {
// 		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid_request_body")
// 	}
// 	if len(req.IDs) == 0 {
// 		return ApiResponse.Error(c, fiber.StatusBadRequest, "no_ids_provided")
// 	}
// 	err := gc.Service.DeleteMultiple(req.IDs)
// 	if err != nil {
// 		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
// 	}
// 	return ApiResponse.Success(c, 200, map[string]interface{}{}, "announcements_deleted_successfully")
// }

// Deletes godoc
// @Summary      delete  announcements
// @Description   delete announcements
// @Tags         Announcements
// @Param        payload body requests.DeleteAnnouncementsRequest true "Delete announcements"
// @Success      200 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /announcements/multiple/delete [post]
func (gc *AnnouncementController) DeleteMultipleAnnouncement(c *fiber.Ctx) error {
	idsParam := c.Query("ids") // e.g. "1,2,3"

	if idsParam == "" {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "no_ids_provided")
	}

	// Convert string → []uint
	idStrs := strings.Split(idsParam, ",")
	var ids []uint

	for _, idStr := range idStrs {
		id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
		if err != nil {
			return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid_id_format")
		}
		ids = append(ids, uint(id))
	}

	if len(ids) == 0 {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "no_valid_ids")
	}
	user, _ := c.Locals("user").(jwt.MapClaims)
	err := gc.Service.DeleteMultiple(ids)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	gc.Service.AuditAction(
		c,
		"mass delete",
		"announcement detail",
		user,
		ids,
		nil,
	)

	return ApiResponse.Success(c, 200, map[string]interface{}{}, "announcements_deleted_successfully")
}

// GetAnnouncementsListForFE godoc
// @Summary      List Announcements for FE
// @Description  Get enabled announcements for a specific broadcast audience and language
// @Tags         Announcements
// @Produce      json
// @Param        broadcast_audience   query     string  true   "Audience to filter announcements (required)"
// @Param        Accept-Language      header    string  false  "Language code (en, km)" default(en)
// @Success      200 {object} ApiResponse.Response "Successful response with announcements"
// @Failure      400 {object} ApiResponse.Response "Missing required query parameter broadcast_audience"
// @Failure      500 {object} ApiResponse.Response "Internal server error"
// @Router       /announcement-list [get]
func (gc *AnnouncementController) GetAnnouncementsListForFE(c *fiber.Ctx) error {
	lang := c.Get("Accept-Language", "en")
	broadcastAudience := c.Query("broadcast_audience") // required
	if broadcastAudience == "" {
		return ApiResponse.JsonError(c, fiber.StatusBadRequest, "broadcast_audience_is_required")
	}
	announcements, err := gc.Service.GetAllForFE(broadcastAudience, lang)
	if err != nil {
		return ApiResponse.JsonError(c, fiber.StatusInternalServerError, err.Error())
	}
	resp := make([]Resources.AnnouncementListFEResource, 0, len(announcements))
	for i := range announcements {
		resp = append(resp, mapAnnouncementFEToResouce(&announcements[i], lang))
	}
	return ApiResponse.JsonSuccess(c, fiber.StatusOK, resp, "announcements_retrieved_successfully")
}

// DeleteAnnouncementDetail godoc
// @Summary      Delete announcementdetails
// @Description  Delete announcementdetails by ID
// @Tags         Announcements
// @Param        id path int true "announcementDetailID"
// @Success      204 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /announcements/detail/delete/{id} [delete]
func (gc *AnnouncementController) DeleteAnnouncementDetail(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	log.Println("hello delete")
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid_id")
	}
	if id < 0 {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid_id")
	}
	oldAnnouncementDetail, _ := gc.Service.GetAnnouncementDetailByID(uint(id))
	announcementnew, err := gc.Service.GetAnnouncementByID(uint(oldAnnouncementDetail.AnnouncementID))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	user, _ := c.Locals("user").(jwt.MapClaims)
	err = gc.Service.DeleteDetail(uint(id))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}

	oldlog := mapAnnouncementDetailLogToResouce(oldAnnouncementDetail, announcementnew)
	gc.Service.AuditAction(
		c,
		"delete",
		"announcement detail",
		user,
		oldlog,
		nil,
	)

	// Return empty success response for consistency
	return ApiResponse.Success(c, 200, map[string]interface{}{}, "announcement_detail_deleted_successfully")
}

// GetAnnouncementsDetail godoc
// @Summary      Get Announcement Detail
// @Description  Retrieve a single announcement by ID
// @Tags         Announcements
// @Produce      json
// @Param        id                   path      int     true   "Announcement ID"
// @Param        Accept-Language      header    string  false  "Language code (en, km)" default(en)
// @Success      200 {object} ApiResponse.Response "Successful response with announcement data"
// @Failure      400 {object} ApiResponse.Response "Invalid ID"
// @Failure      404 {object} ApiResponse.Response "Announcement not found"
// @Failure      500 {object} ApiResponse.Response "Internal server error"
// @Router       /announcements/{id} [get]
func (gc *AnnouncementController) GetAnnouncementsDetail(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id < 0 {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid_id")
	}

	announcement, err := gc.Service.GetAnnouncementsDetail(uint(id))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}

	resp := mapAnnouncementToResouce(announcement)

	// Wrap single record in "records" array to match list structure
	return ApiResponse.Success(c, fiber.StatusOK, resp, "announcement_data_retrieved_successfully")
}
