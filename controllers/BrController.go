package controllers

import (
	"config-service/ApiResponse"
	"config-service/Repositories"
	"config-service/Services"
	"config-service/config"
	"config-service/helper"
	"config-service/requests"
	"config-service/utils"
	"encoding/json"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type BannerController struct {
	Service *Services.BannerService
}

func NewBannerController(db *gorm.DB) *BannerController {
	auditService := &helper.AuditService{
		//Redis: redisConn,
		Channel: config.RabbitChannel,
		Queue:   "audit_logs",
	}
	return &BannerController{
		Service: &Services.BannerService{
			Repo: &Repositories.BannerRepository{
				DB: db,
			},
			Audit: auditService,
		},
	}
}

// CreateBanner godoc
// @Summary      Create a Banner
// @Description  Add new banner with multilingual URLs
// @Tags         Banners
// @Accept       json
// @Produce      json
// @Param        payload body requests.BannerCreateRequest true "Banner Payload"
// @Success      201 {object} Resources.BannerResource
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      422 {object} ApiResponse.ValidationResponse
// @Router       /banners [post]
func (b *BannerController) CreateBanner(c *fiber.Ctx) error {
	req := new(requests.BannerCreateRequest)
	user := helper.ParseStringToUint(c.Get("X-ID"))
	if user <= 0 {
		return ApiResponse.Error(c, fiber.StatusUnauthorized, "user_not_authenticated")
	}

	if err := json.Unmarshal(c.Body(), req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("Controller bad request", err.Error())
		return ApiResponse.ValidationError(c, fiber.StatusBadRequest, "unable_to_parse_request_body", fullMessages[0])
	}

	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("Controller err", fullMessages[0])
		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation_error", fullMessages[0])
	}

	if req.DisplayDuration == "schedule" && (req.StartDate == nil || req.EndDate == nil) {
		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation_error", "start_and_end_date_are_required_for_scheduled_banners")
	}

	banner, err := b.Service.Create(req, user)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	//creator, _ := c.Locals("user").(jwt.MapClaims)
	creator := jwt.MapClaims{
		"email": c.Get("X-User-Email"),
		"role":  c.Get("X-User-Role"),
	}
	newBanner, _ := b.Service.GetByID(uint64(banner.ID))
	b.Service.AuditAction(
		c,
		"create",
		"banner",
		creator,
		nil,
		newBanner,
	)

	return ApiResponse.Success(c, fiber.StatusCreated, banner, "banner_created_successfully")
}

// GetAllBanners godoc
// @Summary      Get all Banners
// @Description  Retrieve list of all banners
// @Tags         Banners
// @Produce      json
// @Success      200 {object} Resources.BannerResource
// @Failure      500 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      422 {object} ApiResponse.ValidationResponse
// @Router       /banners [get]
func (b *BannerController) GetAllBanners(c *fiber.Ctx) error {

	req := new(requests.BannerSearchRequest)
	if err := c.QueryParser(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		return ApiResponse.ValidationError(c, fiber.StatusBadRequest, "unable_to_parse_request_body", fullMessages[0])
	}
	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("err", fullMessages[0])
		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation_error", fullMessages[0])
	}

	banners, err := b.Service.GetAll(c, req)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return ApiResponse.Pagination(c, fiber.StatusOK, banners, "banners_fetched_successfully")
}

// GetBannerByID godoc
// @Summary      Get Banner by ID
// @Description  Retrieve a single banner by ID
// @Tags         Banners
// @Produce      json
// @Param        id path int true "Banner ID"
// @Success      200 {object} Resources.BannerResource
// @Failure      404 {object} ApiResponse.Response
// @Router       /banners/{id} [get]
func (b *BannerController) GetBannerByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id < 0 {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid_id")
	}
	banner, err := b.Service.GetByID(uint64(id))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, "banner_not_found")
	}

	return ApiResponse.Success(c, fiber.StatusOK, banner, "banner_fetched_successfully")
}

// UpdateBanner godoc
// @Summary      Update Banner
// @Description  Update banner details by ID
// @Tags         Banners
// @Accept       json
// @Produce      json
// @Param        id path int true "Banner ID"
// @Param        payload body requests.BannerUpdateRequest true "Banner Update Payload"
// @Success      200 {object} Resources.BannerResource
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      422 {object} ApiResponse.ValidationResponse
// @Router       /banners/{id} [put]
func (b *BannerController) UpdateBanner(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid_id")
	}
	user := helper.ParseStringToUint(c.Get("X-ID"))
	if user <= 0 {
		return ApiResponse.Error(c, fiber.StatusUnauthorized, "user_not_authenticated")
	}
	req := new(requests.BannerUpdateRequest)

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
	oldBanner, _ := b.Service.GetByID(uint64(id))
	_, errNew := b.Service.Update(uint64(id), req, uint64(user))

	if errNew != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, errNew.Error())
	}
	newBanner, _ := b.Service.GetByID(uint64(id))
	//creator, _ := c.Locals("user").(jwt.MapClaims)
	creator := jwt.MapClaims{
		"email": c.Get("X-User-Email"),
		"role":  c.Get("X-User-Role"),
	}
	b.Service.AuditAction(
		c,
		"update",
		"banner",
		creator,
		oldBanner,
		newBanner,
	)

	return ApiResponse.Success(c, fiber.StatusOK, newBanner, "banner_updated_successfully")
}

// DeleteBanner godoc
// @Summary      Delete Banner
// @Description  Delete a banner by ID
// @Tags         Banners
// @Produce      json
// @Param        id path int true "Banner ID"
// @Success      200 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Router       /banners/{id} [delete]
func (b *BannerController) DeleteBanner(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id < 0 {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid _id")
	}
	deleteBanner, _ := b.Service.GetByID(uint64(id))
	if err := b.Service.Delete(uint64(id)); err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	//creator, _ := c.Locals("user").(jwt.MapClaims)
	creator := jwt.MapClaims{
		"email": c.Get("X-User-Email"),
		"role":  c.Get("X-User-Role"),
	}
	b.Service.AuditAction(
		c,
		"delete",
		"banner",
		creator,
		deleteBanner,
		nil,
	)
	return ApiResponse.Success(c, fiber.StatusOK, nil, "banner_deleted_successfully")
}

// MassUpdateBannerStatus godoc
// @Summary      Update  Banners status
// @Description  Mass update Banners status
// @Tags         Banners
// @Param        payload body requests.BannerMassStatusUpdateRequest true "Update Banners status"
// @Success      200 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /status [post]
func (b *BannerController) MassUpdateBannerStatus(c *fiber.Ctx) error {
	req := new(requests.BannerMassStatusUpdateRequest)
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
	oldBanners, err := b.Service.GetByIDs(req.BannerIDs)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}

	count, err := b.Service.MassUpdateStatus(
		req.BannerIDs,
		req.Status,
	)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	newBanners, err := b.Service.GetByIDs(req.BannerIDs)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	creator := jwt.MapClaims{
		"email": c.Get("X-User-Email"),
		"role":  c.Get("X-User-Role"),
	}
	oldItems := make([]map[string]interface{}, 0, len(oldBanners))
	for _, banner := range oldBanners {
		oldItems = append(oldItems, map[string]interface{}{
			"image_name": banner.ImageName,
			"status":     banner.Status,
		})
	}

	newItems := make([]map[string]interface{}, 0, len(newBanners))
	for _, banner := range newBanners {
		newItems = append(newItems, map[string]interface{}{
			"image_name": banner.ImageName,
			"status":     banner.Status,
		})
	}

	oldData := map[string]interface{}{
		"banner": oldItems,
	}

	newData := map[string]interface{}{
		"banner": newItems,
	}
	b.Service.AuditAction(
		c,
		"mass update",
		"banner status",
		creator,
		oldData,
		newData,
	)

	return ApiResponse.Success(c, fiber.StatusOK, strconv.FormatInt(count, 10)+" records_updated", "banners_updated_successfully")
}

// UpdateBannerStatus godoc
// @Summary      Update  Banners status
// @Description  update Banners status
// @Tags         Banners
// @Param        payload body requests.BannerStatusUpdateRequest true "Update Banners status"
// @Success      200 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /status [post]
func (b *BannerController) UpdateBannerStatus(c *fiber.Ctx) error {
	req := new(requests.BannerStatusUpdateRequest)
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

	oldBanner, err := b.Service.GetByID(uint64(req.BannerID))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}

	count, err := b.Service.MassUpdateStatus(
		[]uint{req.BannerID},
		req.Status,
	)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	newBanner, err := b.Service.GetByID(uint64(req.BannerID))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}

	oldData := map[string]interface{}{
		"banner": map[string]interface{}{
			"image_name": oldBanner.ImageName,
			"status":     oldBanner.Status,
		},
	}

	newData := map[string]interface{}{
		"banner": map[string]interface{}{
			"image_name": newBanner.ImageName,
			"status":     newBanner.Status,
		},
	}

	creator := jwt.MapClaims{
		"email": c.Get("X-User-Email"),
		"role":  c.Get("X-User-Role"),
	}

	b.Service.AuditAction(
		c,
		"update",
		"banner status",
		creator,
		oldData,
		newData,
	)

	return ApiResponse.Success(c, fiber.StatusOK, strconv.FormatInt(count, 10)+" records_updated", "banners_updated_successfully")
}

// MassBannerDelete godoc
// @Summary      Delete  Banners
// @Description  Mass delete Banners
// @Tags         Banners
// @Param        payload body requests.BannerMassDeleteRequest true "delete Banners "
// @Success      200 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /mass-delete [post]
func (b *BannerController) MassBannerDelete(c *fiber.Ctx) error {
	req := new(requests.BannerMassDeleteRequest)
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

	count, err := b.Service.MassDelete(
		req.BannerIDs,
	)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	creator := jwt.MapClaims{
		"email": c.Get("X-User-Email"),
		"role":  c.Get("X-User-Role"),
	}
	b.Service.AuditAction(
		c,
		"mass delete",
		"banner",
		creator,
		req.BannerIDs,
		nil,
	)

	return ApiResponse.Success(c, fiber.StatusOK, strconv.FormatInt(count, 10)+" records_deleted", "banner_deleted_successfully")
}

// GetActiveBanners godoc
// @Summary      get active  Banners
// @Description  get active Banners within current time
// @Tags         Banners
// @Param        payload body requests.GetBannerRequest true "get Banners "
// @Success      200 {object} Resources.BannerListResource
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /banner-list [post]
func (b *BannerController) GetActiveBanners(c *fiber.Ctx) error {
	req := new(requests.GetBannerRequest)
	langheader := c.Get("Accept-Language")
	lang := helper.NormalizeLanguage(langheader)

	if err := c.BodyParser(&req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("bad request", err.Error())
		return ApiResponse.JsonValidationError(c, fiber.StatusBadRequest, "unable_to_parse_request_body", fullMessages[0])
	}

	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("update err", fullMessages[0])
		return ApiResponse.JsonValidationError(c, fiber.StatusUnprocessableEntity, "validation_error", fullMessages[0])
	}

	banners, err := b.Service.GetActiveBanners(req.DisplayPlatform, lang)

	if err != nil {
		return ApiResponse.JsonError(c, fiber.StatusNotFound, "banner_not_found")
	}

	return ApiResponse.JsonSuccess(c, fiber.StatusOK, banners, "banners_fetched_successfully")
}
