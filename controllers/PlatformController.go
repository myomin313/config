package controllers

import (
	//"time"

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

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// ------------------------ Controller Struct ------------------------

type PlatformController struct {
	Service *Services.PlatformService
}

func NewPlatformController(db *gorm.DB) *PlatformController {
	auditService := &helper.AuditService{
		//Redis: redisConn,
		Channel: config.RabbitChannel,
		Queue:   "audit_logs",
	}
	return &PlatformController{
		Service: &Services.PlatformService{
			Repo: &Repositories.PlatformRepository{
				DB: db,
			},
			Audit: auditService,
		},
	}
}

// ------------------------ Helper ------------------------

func mapPlatformToResponse(p *models.Platform) Resources.PlatformResource {
	ipAllowlist := p.IPAllowlist

	var apiKey, secretKey string

	if p.APIKey != nil {
		apiKey = p.APIKey.Key
		secretKey = p.APIKey.Secret
	}

	return Resources.PlatformResource{
		ID:          p.ID,
		Name:        p.Name,
		Email:       p.Email,
		Status:      p.Status,
		IPAllowlist: ipAllowlist,
		ApiKey:      apiKey,
		SecretKey:   secretKey,
		CreatedAt:   p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
func mapMerchantCurrencyLogToResponse(c *models.MerchantCurrency) Resources.MerchantCurrencyLogResource {

	return Resources.MerchantCurrencyLogResource{
		ID:           c.ID,
		Currency:     c.Currency,
		MerchantCode: c.MerchantCode,
		Status:       c.Status,
	}
}

// ------------------------ Handlers ------------------------

// CreatePlatform godoc
// @Summary      Create a Platform
// @Description  Create a new platform
// @Tags         Platforms
// @Accept       json
// @Produce      json
// @Param        payload body requests.PlatformCreateRequest true "Platform Data"
// @Success      201 {object} Resources.PlatformResource
// @Failure      400 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /platforms [post]
func (pc *PlatformController) CreatePlatform(c *fiber.Ctx) error {
	req := new(requests.PlatformCreateRequest)
	if err := c.BodyParser(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		return ApiResponse.ValidationError(c, fiber.StatusBadRequest, "unable to parse request body", fullMessages[0])
	}
	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation error", fullMessages[0])
	}
	platform, err := pc.Service.Create(req.Name, req.Email, req.Status, req.IPAllowlist, req.ApiKey, req.SecretKey)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusUnprocessableEntity, err.Error())
	}
	platformData, err := pc.Service.FindByID(platform.ID)
	//log.Println("platforms.APIKey", platforms.APIKey)
	resp := mapPlatformToResponse(platform)
	newresp := mapPlatformToResponse(platformData)
	creator, _ := c.Locals("user").(jwt.MapClaims)
	pc.Service.AuditAction(
		c,
		"create",
		"merchant",
		creator,
		nil,
		newresp,
	)
	return ApiResponse.Success(c, 201, resp, "Platform created successfully")
}

// GetPlatforms godoc
// @Summary      List Platforms
// @Description  Get all platforms
// @Tags         Platforms
// @Produce      json
// @Success      200 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /platforms [get]
func (pc *PlatformController) GetPlatforms(c *fiber.Ctx) error {
	platforms, err := pc.Service.All()
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	resp := make([]Resources.PlatformResource, 0, len(platforms))
	for i := range platforms {
		resp = append(resp, mapPlatformToResponse(&platforms[i]))
	}

	return ApiResponse.SuccessWithPagination(c, 200, len(platforms), resp, "Platforms retrieved successfully")
}

// GetPlatform godoc
// @Summary      Get a single Platform
// @Description  Get platform by ID
// @Tags         Platforms
// @Produce      json
// @Param        id path int true "Platform ID"
// @Success      200 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /platforms/{id} [get]
func (pc *PlatformController) GetPlatform(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "Invalid ID")
	}

	platform, err := pc.Service.FindByID(uint64(id))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, "Platform not found")
	}

	resp := mapPlatformToResponse(platform)
	return ApiResponse.Success(c, 200, resp, "Platform retrieved successfully")
}

// UpdatePlatform godoc
// @Summary      Update Platform
// @Description  Update platform by ID
// @Tags         Platforms
// @Accept       json
// @Produce      json
// @Param        id path int true "Platform ID"
// @Param        payload body requests.PlatformUpdateRequest true "Updated Platform Data"
// @Success      200 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /platforms/{id} [put]
func (pc *PlatformController) UpdatePlatform(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "Invalid ID")
	}

	req := new(requests.PlatformUpdateRequest)
	if err := c.BodyParser(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		return ApiResponse.ValidationError(c, fiber.StatusBadRequest, "unable to parse request body", fullMessages[0])
	}
	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation error", fullMessages[0])
	}
	oldPlatform, err := pc.Service.FindByID(uint64(id))
	oldresp := mapPlatformToResponse(oldPlatform)

	platform, err := pc.Service.Update(uint64(id), req.Name, req.Email, req.Status, req.IPAllowlist, req.ApiKey,
		req.SecretKey)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, err.Error())
	}

	resp := mapPlatformToResponse(platform)
	platformData, err := pc.Service.FindByID(platform.ID)
	//log.Println("platforms.APIKey", platforms.APIKey)
	newresp := mapPlatformToResponse(platformData)

	creator, _ := c.Locals("user").(jwt.MapClaims)
	pc.Service.AuditAction(
		c,
		"update",
		"merchant",
		creator,
		oldresp,
		newresp,
	)
	return ApiResponse.Success(c, 200, resp, "Platform updated successfully")
}

// DeletePlatform godoc
// @Summary      Delete Platform
// @Description  Delete platform by ID
// @Tags         Platforms
// @Param        id path int true "Platform ID"
// @Success      204 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /platforms/{id} [delete]
func (pc *PlatformController) DeletePlatform(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "Invalid ID")
	}
	deletePlatform, err := pc.Service.FindByID(uint64(id))
	oldresp := mapPlatformToResponse(deletePlatform)
	err = pc.Service.Delete(uint64(id))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	creator, _ := c.Locals("user").(jwt.MapClaims)
	pc.Service.AuditAction(
		c,
		"delete",
		"merchant",
		creator,
		oldresp,
		nil,
	)

	// Return empty success response for consistency
	return ApiResponse.Success(c, 200, map[string]interface{}{}, "Platform deleted successfully")
}

// GetMerchantCurrency godoc
// @Summary      List Currecny
// @Description  Get all platforms
// @Tags         Platforms
// @Produce      json
// @Success      200 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /currencies [get]
func (pc *PlatformController) GetMerchantCurrency(c *fiber.Ctx) error {
	currencies, err := pc.Service.GetMerchantCurrency()
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	resp := make([]Resources.CurrencyResource, 0, len(currencies))
	for _, currency := range currencies {
		resp = append(resp, Resources.CurrencyResource{
			ID:       currency.ID,
			Currecny: currency.Currency,
		})
	}

	return ApiResponse.Success(c, 200, resp, "currencies retrieved successfully")
}
func (pc *PlatformController) GetMerchantList(c *fiber.Ctx) error {
	platformID, err := c.ParamsInt("id")
	if err != nil || platformID < 0 {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid_id")
	}
	currencies, err := pc.Service.GetMerchantList(uint64(platformID))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	resp := make([]Resources.MerchantResource, 0, len(currencies))
	for _, currency := range currencies {
		resp = append(resp, Resources.MerchantResource{
			ID:           currency.ID,
			Currecny:     currency.Currency,
			MerchantCode: currency.MerchantCode,
			Name:         currency.Platform.Name,
			PlatformID:   currency.PlatformID,
			Email:        currency.Platform.Email,
			Status:       currency.Status,
		})
	}

	return ApiResponse.Success(c, 200, resp, "merchant retrieved successfully")
}
func (pc *PlatformController) CreateMerchantCurrency(c *fiber.Ctx) error {

	platformID, err := c.ParamsInt("id")
	if err != nil {
		return ApiResponse.Error(
			c,
			fiber.StatusBadRequest,
			"invalid platform id",
		)
	}

	req := new(requests.CreateMerchantCurrencyRequest)

	if err := c.BodyParser(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		return ApiResponse.ValidationError(c, fiber.StatusBadRequest, "unable to parse request body", fullMessages[0])
	}
	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation error", fullMessages[0])
	}

	result, err := pc.Service.CreateMerchantCurrency(
		uint(platformID),
		req.MerchantCode,
		req.Currency,
		req.Status,
	)
	if err != nil {
		return ApiResponse.Error(
			c,
			fiber.StatusUnprocessableEntity,
			err.Error(),
		)
	}
	creator, _ := c.Locals("user").(jwt.MapClaims)
	log.Println("result", result)
	newlog := mapMerchantCurrencyLogToResponse(result)
	log.Println("newlog", newlog)
	pc.Service.AuditAction(
		c,
		"create",
		"merchant currency",
		creator,
		nil,
		newlog,
	)

	return ApiResponse.Success(
		c,
		200,
		result,
		"Merchant currency created successfully",
	)
}
func (pc *PlatformController) UpdateMerchantCurrency(c *fiber.Ctx) error {

	merchantCurrencyID, err := c.ParamsInt("id")
	if err != nil {
		return ApiResponse.Error(
			c,
			fiber.StatusBadRequest,
			"invalid merchant currency id",
		)
	}

	req := new(requests.UpdateMerchantCurrencyRequest)

	if err := c.BodyParser(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		return ApiResponse.ValidationError(
			c,
			fiber.StatusBadRequest,
			"unable to parse request body",
			fullMessages[0],
		)
	}

	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		return ApiResponse.ValidationError(
			c,
			fiber.StatusUnprocessableEntity,
			"validation error",
			fullMessages[0],
		)
	}
	oldMerchantCurrency, err := pc.Service.FindMerchantCurrencyByID(uint(merchantCurrencyID))
	oldlog := mapMerchantCurrencyLogToResponse(oldMerchantCurrency)
	result, err := pc.Service.UpdateMerchantCurrency(
		uint(merchantCurrencyID),
		req.MerchantCode,
		req.Currency,
		req.Status,
	)
	newlog := mapMerchantCurrencyLogToResponse(result)
	if err != nil {
		return ApiResponse.Error(
			c,
			fiber.StatusUnprocessableEntity,
			err.Error(),
		)
	}
	creator, _ := c.Locals("user").(jwt.MapClaims)
	pc.Service.AuditAction(
		c,
		"update",
		"merchant currency",
		creator,
		oldlog,
		newlog,
	)

	return ApiResponse.Success(
		c,
		200,
		result,
		"Merchant currency updated successfully",
	)
}

func (pc *PlatformController) GetAllPlatforms(c *fiber.Ctx) error {
	platforms, err := pc.Service.AllEnable()
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	resp := make([]Resources.PlatformsResource, 0, len(platforms))
	for i := range platforms {
		resp = append(resp, Resources.PlatformsResource{
			ID:   platforms[i].ID,
			Name: platforms[i].Name,
		})
	}

	return ApiResponse.Success(c, 200, resp, "Platforms retrieved successfully")
}
