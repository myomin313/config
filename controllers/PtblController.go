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
	"os"
	"time"

	"config-service/pb/commonpb"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ------------------------ Controller Struct ------------------------

type PlayTableController struct {
	Service         *Services.PlayTableService
	PlayHallService *Services.PlayHallService
	PlayService     *Services.PlayService
	PlayTypeService *Services.PlayTypeService
}

func NewPlayTableController(db *gorm.DB, redisConn *redis.Client) *PlayTableController {
	auditService := &helper.AuditService{
		//Redis: redisConn,
		Channel: config.RabbitChannel,
		Queue:   "audit_logs",
	}
	return &PlayTableController{
		Service: &Services.PlayTableService{
			Repo: &Repositories.PlayTableRepository{
				DB: db,
			},
			PlayRepo: &Repositories.PlayRepository{
				DB: db,
			},
			RedisPublisher: Services.NewRedisPublisher(
				redisConn,
				300*time.Millisecond,
			),
			Audit: auditService,
		},
		PlayHallService: &Services.PlayHallService{
			Repo: &Repositories.PlayHallRepository{
				DB: db,
			},
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

func mapPlayTableToResponse(p *models.GameTable) Resources.GameTableResource {

	var createdAt, updatedAt, mobileIcon, pcIcon string
	if !p.CreatedAt.IsZero() {
		createdAt = p.CreatedAt.Format("2006-01-02 15:04:05")
	}

	if !p.UpdatedAt.IsZero() {
		updatedAt = p.UpdatedAt.Format("2006-01-02 15:04:05")
	}

	if p.MobileIcon != "" {
		mobileIcon = os.Getenv("MINIO_PUBLIC_ENDPOINT") + "/" + p.MobileIcon
	}

	if p.PCIcon != "" {
		pcIcon = os.Getenv("MINIO_PUBLIC_ENDPOINT") + "/" + p.PCIcon
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
	// Declare pointer for UpdatedBy
	var respupdate *commonpb.GetUserResponse
	respupdate = resp

	if p.CreatedByID != p.UpdatedByID {
		respupdate, err = client.GetUser(ctx, &commonpb.GetUserRequest{Id: int64(p.UpdatedByID)})
		if err != nil {
			log.Println("Error calling GetUser:", err)
		}
	}

	log.Println("p.Resolutions:", p.Resolutions)
	resolutions := make(
		[]Resources.GameTableResolutionResource,
		0,
		len(p.Resolutions),
	)

	for _, resolution := range p.Resolutions {
		resolutions = append(
			resolutions,
			Resources.GameTableResolutionResource{
				ID:          resolution.ID,
				GameTableID: resolution.GameTableID,
				Resolution:  resolution.Resolution,
				RoomName:    resolution.RoomName,
			},
		)
	}

	return Resources.GameTableResource{
		ID:              p.ID,
		ENName:          getPlayTableNameByLang(p.Translations, "en"),
		CHName:          getPlayTableNameByLang(p.Translations, "ch"),
		KHName:          getPlayTableNameByLang(p.Translations, "kh"),
		AffiliateOffice: p.AffiliateOffice,
		GameType:        p.GameType,
		GameTypeID:      p.GameTypeID,
		BettingTime:     p.BettingTime,
		SettlementTime:  p.SettlementTime,
		TableNumber:     p.TableNumber,
		Status:          p.Status,
		MobileIcon:      mobileIcon,
		//TableResultType:   p.TableResultType,
		PCIcon:            pcIcon,
		Sort:              p.Sort,
		CreatedBy:         resp.Name,
		UpdatedBy:         respupdate.Name,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		Resolutions:       resolutions,
		Mode:              p.Mode,
		Topic:             p.Topic,
		KafkaOnlineStatus: p.KafkaOnlineStatus,
	}

}

func mapPlayTableToFEResponse(p *models.GameTable, lang string) Resources.GameTableFEResource {

	var mobileIcon, pcIcon string

	if p.MobileIcon != "" {
		mobileIcon = os.Getenv("MINIO_PUBLIC_ENDPOINT") + "/" + p.MobileIcon
	}

	if p.PCIcon != "" {
		pcIcon = os.Getenv("MINIO_PUBLIC_ENDPOINT") + "/" + p.PCIcon
	}

	return Resources.GameTableFEResource{
		ID:              p.ID,
		Name:            getPlayTableNameByLang(p.Translations, lang),
		AffiliateOffice: p.AffiliateOffice,
		GameType:        p.GameType,
		BettingTime:     p.BettingTime,
		SettlementTime:  p.SettlementTime,
		TableNumber:     p.TableNumber,
		Status:          p.Status,
		MobileIcon:      mobileIcon,
		//TableResultType: p.TableResultType,
		PCIcon: pcIcon,
		Sort:   p.Sort,
	}

}
func mapPlayTableLogToResponse(p *models.GameTable, hall *models.GameHall, play *models.Game, playtype *models.GameType) Resources.GameTableLogResource {

	var createdAt, updatedAt, mobileIcon, pcIcon string
	if !p.CreatedAt.IsZero() {
		createdAt = p.CreatedAt.Format("2006-01-02 15:04:05")
	}

	if !p.UpdatedAt.IsZero() {
		updatedAt = p.UpdatedAt.Format("2006-01-02 15:04:05")
	}

	if p.MobileIcon != "" {
		mobileIcon = os.Getenv("MINIO_PUBLIC_ENDPOINT") + "/" + p.MobileIcon
	}

	if p.PCIcon != "" {
		pcIcon = os.Getenv("MINIO_PUBLIC_ENDPOINT") + "/" + p.PCIcon
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
	// Declare pointer for UpdatedBy
	var respupdate *commonpb.GetUserResponse
	respupdate = resp

	if p.CreatedByID != p.UpdatedByID {
		respupdate, err = client.GetUser(ctx, &commonpb.GetUserRequest{Id: int64(p.UpdatedByID)})
		if err != nil {
			log.Println("Error calling GetUser:", err)
		}
	}

	log.Println("p.Resolutions:", p.Resolutions)
	resolutions := make(
		[]Resources.GameTableResolutionLogResource,
		0,
		len(p.Resolutions),
	)

	for _, resolution := range p.Resolutions {
		resolutions = append(
			resolutions,
			Resources.GameTableResolutionLogResource{
				Resolution: resolution.Resolution,
				RoomName:   resolution.RoomName,
			},
		)
	}

	result := Resources.GameTableLogResource{
		ENName: getPlayTableNameByLang(p.Translations, "en"),
		CHName: getPlayTableNameByLang(p.Translations, "ch"),
		KHName: getPlayTableNameByLang(p.Translations, "kh"),
		//AffiliateOffice: p.AffiliateOffice,
		//GameType:        p.GameType,
		//GameTypeID:      p.GameTypeID,
		BettingTime:    p.BettingTime,
		SettlementTime: p.SettlementTime,
		TableNumber:    p.TableNumber,
		Status:         p.Status,
		MobileIcon:     mobileIcon,
		//TableResultType:   p.TableResultType,
		PCIcon:            pcIcon,
		Sort:              p.Sort,
		CreatedBy:         resp.Name,
		UpdatedBy:         respupdate.Name,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		GamePlayType:      playtype.GameTypeName,
		Resolution:        resolutions,
		Mode:              p.Mode,
		Topic:             p.Topic,
		KafkaOnlineStatus: p.KafkaOnlineStatus,
	}

	if hall != nil {
		result.GameHall = &Resources.GameHallLogResource{
			ENName: getNameByLang(hall.Translations, "en"),
			CHName: getNameByLang(hall.Translations, "ch"),
			KHName: getNameByLang(hall.Translations, "kh"),
		}
	}
	if play != nil {
		result.Game = &Resources.GameLogResource{
			ENName: getNameByLang(play.Translations, "en"),
			CHName: getNameByLang(play.Translations, "ch"),
			KHName: getNameByLang(play.Translations, "kh"),
		}
	}
	log.Println("game", play)
	return result

}

func getPlayTableNameByLang(translations []models.GameTranslation, lang string) string {
	for _, t := range translations {
		if t.LanguageCode == lang {
			return t.Content
		}
	}

	return ""
}

// ------------------------ Handlers ------------------------

// Create Game table godoc
// @Summary      Create a Gametable
// @Description  Create a new Gametable
// @Tags         Gametables
// @Accept       json
// @Produce      json
// @Param        payload body requests.GameTableCreateRequest true "GameTable Data"
// @Success      201 {object} Resources.GameTableResource
// @Failure      400 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /gametables [post]
func (pc *PlayTableController) CreatePlaytable(c *fiber.Ctx) error {
	req := new(requests.GameTableCreateRequest)
	user, ok := c.Locals("user").(jwt.MapClaims)

	if !ok || user == nil {
		return ApiResponse.Error(c, fiber.StatusUnauthorized, "User not authenticated")
	}

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
	userID, ok := user["id"].(float64)
	if !ok || userID == 0 {
		return ApiResponse.Error(c, fiber.StatusUnauthorized, "Invalid user ID in token")
	}
	platform, err := pc.Service.Create(req, uint(userID))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, err.Error()) // ✅ return here
	}

	newplatform, _ := pc.Service.FindByID(uint64(platform.ID))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	pc.Service.RedisPublisher.UpdateTableInfo(newplatform)
	resp := mapPlayTableToResponse(newplatform)

	updateplayHall, err := pc.PlayHallService.FindByID(uint64(newplatform.AffiliateOffice))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	updateplay, err := pc.PlayService.FindByID(uint64(newplatform.GameType))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	updateplaytype, err := pc.PlayTypeService.FindByID(uint64(newplatform.GameTypeID))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	logresp := mapPlayTableLogToResponse(newplatform, updateplayHall, updateplay, updateplaytype)

	pc.Service.AuditAction(
		c,
		"create",
		"game table",
		user,
		nil,
		logresp,
	)
	return ApiResponse.Success(c, fiber.StatusCreated, resp, "Game Table created successfully")
}

// GetPlatforms godoc
// @Summary      List game tables
// @Description  Get all game tables
// @Tags         Gametables
// @Produce      json
// @Success      200 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /gametables [get]
func (pc *PlayTableController) GetPlayTables(c *fiber.Ctx) error {
	req := new(requests.GameTableSearchRequest)
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

	platforms, err := pc.Service.All(req, lang)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	resp := make([]Resources.GameTableResource, 0, len(platforms))
	for i := range platforms {
		resp = append(resp, mapPlayTableToResponse(&platforms[i]))
	}

	return ApiResponse.SuccessWithPagination(c, 200, len(platforms), resp, "Game Tables retrieved successfully")
}

// GetPlatform godoc
// @Summary      Get a single gametables
// @Description  Get gametable by ID
// @Tags         Gametables
// @Produce      json
// @Param        id path int true "Game Table ID"
// @Success      200 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /gametables/{id} [get]
func (pc *PlayTableController) GetPlayTable(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "Invalid ID")
	}

	platform, err := pc.Service.FindByID(uint64(id))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, "Record not found in Game Table")
		//return ApiResponse.Error(c, fiber.StatusOK, "Record not found in Game Table")
	}

	resp := mapPlayTableToResponse(platform)
	return ApiResponse.Success(c, 200, resp, "Game Table retrieved successfully")
}

// UpdatePlatform godoc
// @Summary      Update gametables
// @Description  Update gametables by ID
// @Tags         Gametables
// @Accept       json
// @Produce      json
// @Param        id path int true "Game Table ID"
// @Param        payload body requests.GameTableUpdateRequest true "Updated gametable Data"
// @Success      200 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /gametables/{id} [put]
func (pc *PlayTableController) UpdatePlayTable(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "Invalid ID")
	}

	req := new(requests.GameTableUpdateRequest)
	if err := c.BodyParser(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		return ApiResponse.ValidationError(c, fiber.StatusBadRequest, "unable to parse request body", fullMessages[0])
	}
	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation error", fullMessages[0])
	}

	user, ok := c.Locals("user").(jwt.MapClaims)
	if !ok || user == nil {
		return ApiResponse.Error(c, fiber.StatusUnauthorized, "User not authenticated")
	}
	userID, ok := user["id"].(float64)
	if !ok || userID == 0 {
		return ApiResponse.Error(c, fiber.StatusUnauthorized, "Invalid user ID in token")
	}
	playtable, err := pc.Service.FindByID(uint64(id))
	//oldresp := mapGameTableLogToResponse(gametable)
	// gameHall, err := pc.GameHallService.FindByID(uint64(gametable.AffiliateOffice))
	// if err != nil {
	// 	return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	// }
	log.Println("gametable.AffiliateOffice", playtable.AffiliateOffice)
	playHall, err := pc.PlayHallService.FindByID(uint64(playtable.AffiliateOffice))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}

	//	oldresp := mapGameTableLogToResponse(gametable, gameHall,)

	play, err := pc.PlayService.FindByID(uint64(playtable.GameType))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}

	playtype, err := pc.PlayTypeService.FindByID(uint64(playtable.GameTypeID))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	log.Println("gametype", playtype)
	oldresp := mapPlayTableLogToResponse(playtable, playHall, play, playtype)

	_, errNew := pc.Service.Update(uint64(id), req, uint(userID))

	if errNew != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, errNew.Error())
	}
	newplatform, _ := pc.Service.FindByID(uint64(id))

	pc.Service.RedisPublisher.UpdateTableInfo(newplatform)
	resp := mapPlayTableToResponse(newplatform)
	log.Println("newplatform.AffiliateOffice", newplatform.AffiliateOffice)
	updateplayHall, err := pc.PlayHallService.FindByID(uint64(newplatform.AffiliateOffice))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	updateplay, err := pc.PlayService.FindByID(uint64(newplatform.GameType))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	updateplaytype, err := pc.PlayTypeService.FindByID(uint64(newplatform.GameTypeID))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	log.Println("gametype", playtype)
	logresp := mapPlayTableLogToResponse(newplatform, updateplayHall, updateplay, updateplaytype)

	//	oldData, newData := utils.StructChanges(oldresp, logresp)
	pc.Service.AuditAction(
		c,
		"update",
		"game table",
		user,
		oldresp,
		logresp,
	)
	return ApiResponse.Success(c, 200, resp, "Game Table updated successfully")
}

// DeleteGameTable godoc
// @Summary      Delete gametables
// @Description  Delete gametables by ID
// @Tags         Gametables
// @Param        id path int true "gametablesID"
// @Success      204 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /gametables/{id} [delete]
func (pc *PlayTableController) DeletePlayTable(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "Invalid ID")
	}
	playtable, err := pc.Service.FindByID(uint64(id))
	deleteplayHall, err := pc.PlayHallService.FindByID(uint64(playtable.AffiliateOffice))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	deleteplay, err := pc.PlayService.FindByID(uint64(playtable.GameType))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	deleteplaytype, err := pc.PlayTypeService.FindByID(uint64(playtable.GameTypeID))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	oldresp := mapPlayTableLogToResponse(playtable, deleteplayHall, deleteplay, deleteplaytype)
	err = pc.Service.Delete(uint64(id))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}

	pc.Service.AuditAction(
		c,
		"delete",
		"game table",
		c.Locals("user").(jwt.MapClaims),
		oldresp,
		nil,
	)

	// Return empty success response for consistency
	return ApiResponse.Success(c, 200, map[string]interface{}{}, "Game Table deleted successfully")
}

// MassUpdateGameTableStatus godoc
// @Summary      Update  gametables status
// @Description  Mass update gametables status
// @Tags         Gametables
// @Param        payload body requests.GameTableMassStatusUpdateRequest true "Update gametables status"
// @Success      200 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      404 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /status/{ids} [post]
func (pc *PlayTableController) MassUpdatePlayTableStatus(c *fiber.Ctx) error {
	req := new(requests.GameTableMassStatusUpdateRequest)
	if err := c.BodyParser(&req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("bad request", err.Error())
		return ApiResponse.ValidationError(c, fiber.StatusBadRequest, "unable to parse request body", fullMessages[0])
	}

	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("update err", fullMessages[0])
		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation error", fullMessages[0])
	}
	oldtablestatus, err := helper.GetTablesByIDs(req.TableIDs, 0, 0)

	err = pc.Service.MassUpdateStatus(
		req.TableIDs,
		req.Status,
	)
	tables, err := helper.GetTablesByIDs(req.TableIDs, 0, 0)
	log.Println("MassUpdateGameTableStatus tables:", tables)
	tableNumbers := make([]string, 0)

	for _, table := range tables {
		tableNumbers = append(tableNumbers, table.TableNumber)
	}

	oldTables := make([]map[string]interface{}, 0)

	for _, table := range oldtablestatus {
		oldTables = append(oldTables, map[string]interface{}{
			"table_number": table.TableNumber,
			"status":       table.Status,
		})
	}

	oldData := map[string]interface{}{
		"game_table": oldTables,
	}

	newTables := make([]map[string]interface{}, 0)

	for _, table := range tables {
		newTables = append(newTables, map[string]interface{}{
			"table_number": table.TableNumber,
			"status":       req.Status,
		})
	}

	newData := map[string]interface{}{
		"game_table": newTables,
	}

	pc.Service.AuditAction(
		c,
		"mass update",
		"game table",
		c.Locals("user").(jwt.MapClaims),
		oldData,
		newData,
	)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}

	return ApiResponse.Success(c, 200, len(req.TableIDs), "Game Table updated successfully")
}

// GetPlatforms godoc
// @Summary      List game tables
// @Description  Get all game tables
// @Tags         Gametables
// @Produce      json
// @Success      200 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /gametable-lists [get]
func (pc *PlayTableController) GetPlayTableForFE(c *fiber.Ctx) error {

	req := new(requests.GameTableFESearchRequest)
	lang := c.Get("Accept-Language", "en")
	if err := c.QueryParser(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		return ApiResponse.ValidationError(c, fiber.StatusBadRequest, "unable to parse request body", fullMessages[0])
	}
	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("err", fullMessages[0])
		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation error", fullMessages[0])
	}

	playtables, err := pc.Service.AllForFE(req, lang)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	resp := make([]Resources.GameTableFEResource, 0, len(playtables))
	for i := range playtables {
		resp = append(resp, mapPlayTableToFEResponse(&playtables[i], lang))
	}
	return ApiResponse.Success(c, 200, resp, "Game Tables retrieved successfully")
}
