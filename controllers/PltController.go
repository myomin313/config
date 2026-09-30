package controllers

import (
	"config-service/ApiResponse"
	"config-service/Repositories"
	Resources "config-service/Resources"
	"config-service/Services"
	"config-service/config"
	"config-service/helper"
	"config-service/models"
	"config-service/requests"
	"config-service/utils"
	"context"
	"log"
	"time"

	"config-service/pb/commonpb"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ------------------------ Controller Struct ------------------------

type PlayRedLimitController struct {
	Service         *Services.PlayRedLimitService
	PlayService     *Services.PlayService
	PlayTypeService *Services.PlayTypeService
}

func NewPlayRedLimitController(db *gorm.DB, redisConn *redis.Client) *PlayRedLimitController {
	auditService := &helper.AuditService{
		//Redis: redisConn,
		Channel: config.RabbitChannel,
		Queue:   "audit_logs",
	}
	return &PlayRedLimitController{
		Service: &Services.PlayRedLimitService{
			Repo: &Repositories.PlayRedLimitRepository{
				DB: db,
			},
			PlayRepo: &Repositories.PlayRepository{
				DB: db,
			},
			Audit: auditService,
		},
		PlayService: &Services.PlayService{
			Repo: &Repositories.PlayRepository{
				DB: db,
			},
		},
		PlayTypeService: &Services.PlayTypeService{
			Repo: &Repositories.PlayTypeRepository{
				DB: db,
			},
		},
	}
}

// ------------------------ Helper ------------------------

func mapPlayRedLimitTableToResponse(p *models.GameRedLimit) Resources.GameRedLimitResource {

	var createdAt, updatedAt string
	if !p.CreatedAt.IsZero() {
		createdAt = p.CreatedAt.Format("2006-01-02 15:04:05")
	}

	if !p.UpdatedAt.IsZero() {
		updatedAt = p.UpdatedAt.Format("2006-01-02 15:04:05")
	}

	// Prepare game_red_limit only if there are limits
	var playRedLimit map[string]Resources.RedLimit
	if len(p.RedLimitTable) > 0 {
		playRedLimit = make(map[string]Resources.RedLimit)
		for _, limit := range p.RedLimitTable {
			playRedLimit[limit.GameRedLimitKey] = Resources.RedLimit{
				MinLimit: limit.MinLimit,
				MaxLimit: limit.MaxLimit,
			}
		}
	} else {
		playRedLimit = nil // do not include in JSON
	}

	conn := config.GrpcClient()
	defer conn.Close()

	client := commonpb.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Call GetUser
	resp, err := client.GetUser(ctx, &commonpb.GetUserRequest{Id: int64(p.CreatedByID)})
	if err != nil {
		log.Println("Error calling GetUser:", err)
	}
	respupdate, err := client.GetUser(ctx, &commonpb.GetUserRequest{Id: int64(p.UpdatedByID)})
	if err != nil {
		log.Println("Error calling GetUser:", err)
	}

	return Resources.GameRedLimitResource{
		ID:                   p.ID,
		EnName:               getNameByLang(p.Translations, "en"),
		ChName:               getNameByLang(p.Translations, "ch"),
		KhName:               getNameByLang(p.Translations, "kh"),
		GameType:             p.GameType,
		GameTypeCurrency:     p.GameTypeCurrency,
		GamePlayType:         p.GamePlayType,
		GamePlayTypeCurrency: p.GamePlayTypeCurrency,
		Status:               p.Status,
		Sort:                 p.Sort,
		CreatedBy:            resp.Name,
		UpdatedBy:            respupdate.Name,
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
		GameRedLimit:         playRedLimit,
	}

}

func mapPlayRedLimitTableLogToResponse(p *models.GameRedLimit, play *models.Game, playtype *models.GameType) Resources.GameRedLimitLogResource {

	var createdAt, updatedAt string
	if !p.CreatedAt.IsZero() {
		createdAt = p.CreatedAt.Format("2006-01-02 15:04:05")
	}

	if !p.UpdatedAt.IsZero() {
		updatedAt = p.UpdatedAt.Format("2006-01-02 15:04:05")
	}

	// Prepare game_red_limit only if there are limits
	var playRedLimit map[string]Resources.RedLimit
	if len(p.RedLimitTable) > 0 {
		playRedLimit = make(map[string]Resources.RedLimit)
		for _, limit := range p.RedLimitTable {
			playRedLimit[limit.GameRedLimitKey] = Resources.RedLimit{
				MinLimit: limit.MinLimit,
				MaxLimit: limit.MaxLimit,
			}
		}
	} else {
		playRedLimit = nil // do not include in JSON
	}

	conn := config.GrpcClient()
	defer conn.Close()

	client := commonpb.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Call GetUser
	resp, err := client.GetUser(ctx, &commonpb.GetUserRequest{Id: int64(p.CreatedByID)})
	if err != nil {
		log.Println("Error calling GetUser:", err)
	}
	respupdate, err := client.GetUser(ctx, &commonpb.GetUserRequest{Id: int64(p.UpdatedByID)})
	if err != nil {
		log.Println("Error calling GetUser:", err)
	}

	result := Resources.GameRedLimitLogResource{
		ID:               p.ID,
		EnName:           getNameByLang(p.Translations, "en"),
		ChName:           getNameByLang(p.Translations, "ch"),
		KhName:           getNameByLang(p.Translations, "kh"),
		GamePlayType:     playtype.GameTypeName,
		Status:           p.Status,
		Sort:             p.Sort,
		CreatedBy:        resp.Name,
		UpdatedBy:        respupdate.Name,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		GameRedLimit:     playRedLimit,
		GameTypeCurrency: GetCurrencyName(p.GameTypeCurrency),
	}
	if play != nil {
		result.Game = &Resources.GameLogResource{
			ENName: getNameByLang(play.Translations, "en"),
			CHName: getNameByLang(play.Translations, "ch"),
			KHName: getNameByLang(play.Translations, "kh"),
		}
	}
	return result

}
func GetCurrencyName(id uint) string {
	switch id {
	case 1:
		return "USD"
	case 2:
		return "KHR"
	case 3:
		return "CNY"
	case 4:
		return "THB"
	default:
		return "Unknown"
	}
}

// ------------------------ Handlers ------------------------

// Create Game table godoc
// @Summary      Create a RedLimit
// @Description  Create a new RedLimit
// @Tags         RedLimits
// @Accept       json
// @Produce      json
// @Param        payload body requests.RedLimitCreateRequest true "RedLimit Data"
// @Success      201 {object} Resources.GameRedLimitResource
// @Failure      400 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /red-limit [post]
func (pc *PlayRedLimitController) CreateRedLimit(c *fiber.Ctx) error {
	req := new(requests.RedLimitCreateRequest)
	user, _ := c.Locals("user").(jwt.MapClaims)
	if err := c.BodyParser(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("Controller bad request", err.Error())
		return ApiResponse.ValidationError(c, fiber.StatusBadRequest, "unable to parse request body", fullMessages[0])
	}
	//Validate if needed
	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("Controller err", fullMessages[0])
		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation error", fullMessages[0])
	}

	platform, err := pc.Service.CreateRedLimit(req, uint(user["id"].(float64)))
	// lang := c.Get("X-User-Lang")
	newplatform, _ := pc.Service.FindRedLimitByID(uint64(platform.ID))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	resp := mapPlayRedLimitTableToResponse(newplatform)
	pc.Service.AuditAction(
		c,
		"create",
		"Limit Configuration",
		user,
		nil,
		resp,
	)
	return ApiResponse.Success(c, fiber.StatusCreated, resp, "Red Limit created successfully")
}

// GetPlatforms godoc
// @Summary      List red limits
// @Description  Get all red limits records
// @Tags         RedLimits
// @Produce      json
// @Success      200 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /red-limit [get]
func (pc *PlayRedLimitController) GetRedLimits(c *fiber.Ctx) error {
	req := new(requests.RedLimitSearchRequest)
	lang := c.Get("X-User-Lang")
	if err := c.QueryParser(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		return ApiResponse.ValidationError(c, fiber.StatusBadRequest, "unable to parse request body", fullMessages[0])
	}
	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("err", fullMessages[0])
		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation error", fullMessages[0])
	}
	platforms, err := pc.Service.AllRedLimits(req, lang)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	resp := make([]Resources.GameRedLimitResource, 0, len(platforms))
	for i := range platforms {
		resp = append(resp, mapPlayRedLimitTableToResponse(&platforms[i]))
	}

	return ApiResponse.SuccessWithPagination(c, 200, len(platforms), resp, "Red Limits retrieved successfully")
}

// GetPlatform godoc
// @Summary      Get a single RedLimits
// @Description  Get RedLimit by ID
// @Tags         RedLimits
// @Produce      json
// @Param        id path int true "Red limit ID"
// @Success      200 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /red-limit/{id} [get]
func (pc *PlayRedLimitController) GetRedLimit(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	// lang := c.Get("X-User-Lang")
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "Invalid ID")
	}

	platform, err := pc.Service.FindRedLimitByID(uint64(id))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, "Record not found")
	}

	resp := mapPlayRedLimitTableToResponse(platform)
	return ApiResponse.Success(c, 200, resp, "Red limit retrieved successfully")
}

// UpdatePlatform godoc
// @Summary      Update RedLimits
// @Description  Update RedLimits by ID
// @Tags         RedLimits
// @Accept       json
// @Produce      json
// @Param        id path int true "Red Limit ID"
// @Param        payload body requests.RedLimitUpdateRequest true "Updated RedLimit Data"
// @Success      200 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /red-limit/{id} [put]
func (pc *PlayRedLimitController) UpdateRedLimit(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "Invalid ID")
	}

	req := new(requests.RedLimitUpdateRequest)
	if err := c.BodyParser(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("bad request", err.Error())
		return ApiResponse.ValidationError(c, fiber.StatusBadRequest, "unable to parse request body", fullMessages[0])
	}

	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("update err", fullMessages[0])
		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation error", fullMessages[0])
	}
	oldplatform, _ := pc.Service.FindRedLimitByID(uint64(id))
	play, err := pc.PlayService.FindByID(uint64(oldplatform.GameType))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	playtype, err := pc.PlayTypeService.FindByID(uint64(oldplatform.GamePlayType))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	oldresp := mapPlayRedLimitTableLogToResponse(oldplatform, play, playtype)
	user, _ := c.Locals("user").(jwt.MapClaims)
	// lang := c.Get("X-User-Lang")
	_, errNew := pc.Service.UpdateRedLimit(uint64(id), req, uint(user["id"].(float64)), c.Get("X-User-Lang"))

	if errNew != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, errNew.Error())
	}
	newplatform, _ := pc.Service.FindRedLimitByID(uint64(id))
	log.Println("newplatform", newplatform)
	resp := mapPlayRedLimitTableToResponse(newplatform)
	updateplay, err := pc.PlayService.FindByID(uint64(newplatform.GameType))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	updateplaytype, err := pc.PlayTypeService.FindByID(uint64(newplatform.GamePlayType))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	newresp := mapPlayRedLimitTableLogToResponse(newplatform, updateplay, updateplaytype)
	pc.Service.AuditAction(
		c,
		"update",
		"Limit Configuration",
		user,
		oldresp,
		newresp,
	)
	return ApiResponse.Success(c, 200, resp, "Red limit updated successfully")
}

// DeleteRedLimit godoc
// @Summary      Delete RedLimits
// @Description  Delete RedLimits by ID
// @Tags         RedLimits
// @Param        id path int true "RedLimitsID"
// @Success      204 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /red-limit/{id} [delete]
func (pc *PlayRedLimitController) DeleteRedLimit(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "Invalid ID")
	}

	oldLog, errNew := pc.Service.FindRedLimitByID(uint64(id))

	if errNew != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, "Record not found")
	}
	err = pc.Service.DeleteRedLimit(uint64(id))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	resp := mapPlayRedLimitTableToResponse(oldLog)
	pc.Service.AuditAction(
		c,
		"delete",
		"Limit Configuration",
		c.Locals("user").(jwt.MapClaims),
		resp,
		nil,
	)

	// Return empty success response for consistency
	return ApiResponse.Success(c, 200, map[string]interface{}{}, "Red limit deleted successfully")
}
