package controllers

import (
	"config-service/ApiResponse"
	"config-service/Repositories"
	"config-service/Resources"
	"config-service/Services"
	"config-service/config"
	"config-service/helper"
	"config-service/models"
	"config-service/pb/commonpb"
	"config-service/requests"
	"config-service/utils"
	"context"
	"log"
	"math"
	"time"

	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PlayHallController struct {
	Service *Services.PlayHallService
}

func NewPlayHallController(db *gorm.DB, redisConn *redis.Client) *PlayHallController {

	auditService := &helper.AuditService{
		//Redis: redisConn,
		Channel: config.RabbitChannel,
		Queue:   "audit_logs",
	}

	return &PlayHallController{
		Service: &Services.PlayHallService{
			Repo: &Repositories.PlayHallRepository{
				DB: db,
			},
			PlayRepo: &Repositories.PlayRepository{
				DB: db,
			},
			Audit: auditService, // ✅ THIS IS REQUIRED
		},
	}
}
func mapPlayHallToResouce(playhall *models.GameHall) Resources.GameHallResource {

	conn := config.GrpcClient()
	defer conn.Close()

	client := commonpb.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Call GetUser
	resp, err := client.GetUser(ctx, &commonpb.GetUserRequest{Id: int64(playhall.CreatedByID)})
	if err != nil {
		log.Println("Error calling GetUser:", err)
	}
	respupdate, err := client.GetUser(ctx, &commonpb.GetUserRequest{Id: int64(playhall.UpdatedByID)})
	if err != nil {
		log.Println("Error calling GetUser:", err)
	}
	formattedID := fmt.Sprintf("GH-%08d", playhall.ID)
	return Resources.GameHallResource{
		ID:          playhall.ID,
		FormattedID: formattedID,
		Status:      playhall.Status,
		SortOrder:   int(playhall.SortOrder),
		CreatedBy:   resp.Name,
		UpdatedBy:   respupdate.Name,
		CreatedAt:   playhall.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   playhall.UpdatedAt.Format("2006-01-02 15:04:05"),
		EnName:      getNameByLang(playhall.Translations, "en"),
		ZhName:      getNameByLang(playhall.Translations, "ch"),
		KmName:      getNameByLang(playhall.Translations, "kh"),
	}

}

func mapPlayHallListResource(playhall *models.GameHall) Resources.GameHallListResource {

	return Resources.GameHallListResource{
		ID: playhall.ID,
		//Name:          getNameByLang(gamehall.Translations, lang),
		EnName: getNameByLang(playhall.Translations, "en"),
		ZhName: getNameByLang(playhall.Translations, "ch"),
		KmName: getNameByLang(playhall.Translations, "kh"),
	}
}

func mapPlayHallListFEResource(playhall *models.GameHall, lang string) Resources.GameHallListFEResource {

	return Resources.GameHallListFEResource{
		ID:        playhall.ID,
		Name:      getNameByLang(playhall.Translations, lang),
		SortOrder: playhall.SortOrder,
		Status:    playhall.Status,
		//EnName: getNameByLang(gamehall.Translations, "en"),
		//ZhName: getNameByLang(gamehall.Translations, "ch"),
		//KmName: getNameByLang(gamehall.Translations, "kh"),
	}
}

// CreateGameHall godoc
// @Summary      Create a GameHall
// @Description  Create a new gamehall
// @Tags         GameHalls
// @Accept       json
// @Produce      json
// @Param        payload body requests.GameHallCreateRequest true "GameHall Data"
// @Success      201 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /gamehalls [post]
func (gc *PlayHallController) CreatePlayHall(c *fiber.Ctx) error {
	req := new(requests.GameHallCreateRequest)

	user, _ := c.Locals("user").(jwt.MapClaims)
	log.Println("user:", user)
	if err := c.BodyParser(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		return ApiResponse.ValidationError(c, fiber.StatusBadRequest, "unable to parse request body", fullMessages[0])
	}

	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("err", fullMessages[0])
		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation error", fullMessages[0])
	}

	//game, err := gc.Service.CreateGame(*req)
	playhall, err := gc.Service.CreatePlayHall(*req, uint(user["id"].(float64)))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, err.Error())
	}
	newLog := map[string]any{
		"status":     playhall.Status,
		"sort_order": playhall.SortOrder,
		"en_name":    getNameByLang(playhall.Translations, "en"),
		"ch_name":    getNameByLang(playhall.Translations, "ch"),
		"kh_name":    getNameByLang(playhall.Translations, "kh"),
	}
	gc.Service.AuditAction(
		c,
		"create",
		"game hall",
		user,
		nil,
		newLog,
	)
	//return ApiResponse.Success(c, fiber.StatusCreated, game, "Game created successfully")
	resp := mapPlayHallToResouce(playhall)
	return ApiResponse.Success(c, fiber.StatusCreated, resp, "Game Hall created successfully")
}

// GetGameHalls godoc
// @Summary      List GameHalls
// @Description  Get all gamehalls
// @Tags         GameHalls
// @Produce      json
// @Success      200 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /gamehalls [get]
func (gc *PlayHallController) GetPlayHalls(c *fiber.Ctx) error {
	// Get all games from service

	//lang := c.Get("Accept-Language", "en")
	playhalls, err := gc.Service.All()
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	// Map games to resources
	resp := make([]Resources.GameHallResource, 0, len(playhalls))
	for i := range playhalls {
		resp = append(resp, mapPlayHallToResouce(&playhalls[i]))
	}

	// Return plain JSON without pagination
	return ApiResponse.SuccessWithPagination(c, 200, len(playhalls), resp, "Game Hall retrieved successfully")
}

// DeleteGameHall godoc
// @Summary      Delete gamehalls
// @Description  Delete gamehalls by ID
// @Tags         GameHalls
// @Param        id path int true "gamehallsID"
// @Success      204 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /gamehalls/{id} [delete]
func (pc *PlayHallController) DeletePlayHall(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "Invalid ID")
	}
	if id < 0 {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "Invalid ID")
	}
	playhallold, err := pc.Service.FindByID(uint64(id))
	oldLog := map[string]any{
		"status":     playhallold.Status,
		"sort_order": playhallold.SortOrder,
		"en_name":    getNameByLang(playhallold.Translations, "en"),
		"ch_name":    getNameByLang(playhallold.Translations, "ch"),
		"kh_name":    getNameByLang(playhallold.Translations, "kh"),
	}
	err = pc.Service.Delete(uint64(id))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	pc.Service.AuditAction(
		c,
		"delete",
		"game hall",
		c.Locals("user").(jwt.MapClaims),
		oldLog,
		nil,
	)

	// Return empty success response for consistency
	return ApiResponse.Success(c, 200, map[string]interface{}{}, "Game Hall deleted successfully")
}

// UpdateGameHall godoc
// @Summary      Update gamehalls
// @Description  Update gamehalls by ID
// @Tags         GameHall
// @Accept       json
// @Produce      json
// @Param        id path int true "Game Hall ID"
// @Param        payload body requests.GameHallUpdateRequest true "Updated gamehall Data"
// @Success      200 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /gamehalls/{id} [put]
func (pc *PlayHallController) UpdatePlayHall(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "Invalid ID")
	}

	req := new(requests.GameHallUpdateRequest)
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

	user, _ := c.Locals("user").(jwt.MapClaims)
	if id < 0 {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "Invalid ID")
	}
	userID := user["id"].(float64)
	if userID < 0 || userID > float64(math.MaxUint) {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "Invalid user Id")
	}
	playhallold, err := pc.Service.Repo.FindByID(uint64(id))
	oldLog := map[string]any{
		"status":     playhallold.Status,
		"sort_order": playhallold.SortOrder,
		"en_name":    getNameByLang(playhallold.Translations, "en"),
		"ch_name":    getNameByLang(playhallold.Translations, "ch"),
		"kh_name":    getNameByLang(playhallold.Translations, "kh"),
	}

	playhalls, errNew := pc.Service.Update(uint64(id), *req, uint(userID))

	if errNew != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, errNew.Error())
	}

	newLog := map[string]any{
		"status":     playhalls.Status,
		"sort_order": playhalls.SortOrder,
		"en_name":    getNameByLang(playhalls.Translations, "en"),
		"ch_name":    getNameByLang(playhalls.Translations, "ch"),
		"kh_name":    getNameByLang(playhalls.Translations, "kh"),
	}
	//oldData, newData := utils.StructChanges(oldLog, newLog)
	pc.Service.AuditAction(
		c,
		"update",
		"game hall",
		user,
		oldLog,
		newLog,
	)

	resp := mapPlayHallToResouce(playhalls)
	return ApiResponse.Success(c, 200, resp, "Game Hall updated successfully")
}

// GetGameHalls godoc
// @Summary      List GameHalls
// @Description  Get all gamehalls
// @Tags         GameHalls
// @Produce      json
// @Success      200 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /gamehalls/list [get]
func (gc *PlayHallController) GetPlayHallsList(c *fiber.Ctx) error {
	// Get all games from service

	//lang := c.Get("Accept-Language", "en")
	playhalls, err := gc.Service.AllEnable()
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	// Map games to resources
	resp := make([]Resources.GameHallListResource, 0, len(playhalls))
	for i := range playhalls {
		resp = append(resp, mapPlayHallListResource(&playhalls[i]))
	}

	// Return plain JSON without pagination
	return ApiResponse.SuccessWithPagination(c, 200, len(playhalls), resp, "Game Hall retrieved successfully")
}

// GetGameHallsForFE godoc
// @Summary      List GameHalls
// @Description  Get all gamehalls
// @Tags         GameHalls
// @Produce      json
// @Success      200 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /gamehall-list [get]
func (gc *PlayHallController) GetPlayHallsListForFE(c *fiber.Ctx) error {
	// Get all games from service

	lang := c.Get("Accept-Language", "en")
	playhalls, err := gc.Service.AllEnable()
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	// Map games to resources
	resp := make([]Resources.GameHallListFEResource, 0, len(playhalls))
	for i := range playhalls {
		resp = append(resp, mapPlayHallListFEResource(&playhalls[i], lang))
	}

	// Return plain JSON without pagination
	return ApiResponse.Success(c, 200, resp, "Game Hall retrieved successfully")
}
