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
	"math"
	"os"
	"sort"
	"strings"
	"time"

	"config-service/pb/commonpb"

	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ------------------------ Controller Struct ------------------------
type PlayController struct {
	Service *Services.PlayService
}

func NewPlayController(db *gorm.DB, redisConn *redis.Client) *PlayController {

	auditService := &helper.AuditService{
		//Redis: redisConn,
		Channel: config.RabbitChannel,
		Queue:   "audit_logs",
	}
	return &PlayController{
		Service: &Services.PlayService{
			Repo: &Repositories.PlayRepository{
				DB: db,
			},
			Audit: auditService,
		},
	}
}

func mapPlayToResource(play *models.Game) Resources.GameResource {
	// Map GameTypeData grouped by game_type_id
	playTypeDataMap := make(map[uint64][]Resources.GameTypeDataResource)

	for _, gtd := range play.GameTypeData {
		master := gtd.MasterGameType
		playTypeID := master.GameTypeID // <-- group by this

		image := ""
		if gtd.Image != "" {
			image = os.Getenv("MINIO_PUBLIC_ENDPOINT") + "/" + gtd.Image
		}

		resource := Resources.GameTypeDataResource{
			ID:               gtd.ID,
			Odds:             gtd.Odds,
			NameKm:           gtd.NameKm,
			MasterGameTypeID: gtd.MasterGameTypeID,
			Image:            gtd.Image,
			ImageUrl:         image,
			KeyName:          master.KeyName,
			NameEn:           getPlayTypeNameByLang(master.GameTypeTranslation, "en"),
			NameZh:           getPlayTypeNameByLang(master.GameTypeTranslation, "ch"),
		}

		playTypeDataMap[playTypeID] = append(playTypeDataMap[playTypeID], resource)
	}

	// Convert map to slice for response
	groupedPlayTypeData := make([]Resources.GameTypeGroupResource, 0, len(playTypeDataMap))
	for playTypeID, items := range playTypeDataMap {
		playTypeName := ""

		// Find the actual GameTypeName from game.GameTypes
		for _, gt := range play.GameTypes {
			if gt.ID == playTypeID {
				playTypeName = gt.GameTypeName
				break
			}
		}
		sortPlayTypeDataResources(items)
		groupedPlayTypeData = append(groupedPlayTypeData, Resources.GameTypeGroupResource{
			GameTypeID:   playTypeID,
			GameTypeName: playTypeName, // now filled
			GameTypeData: items,
		})
	}

	sort.Slice(groupedPlayTypeData, func(i, j int) bool {
		return groupedPlayTypeData[i].GameTypeID < groupedPlayTypeData[j].GameTypeID
	})

	// --- GRPC calls for CreatedBy / UpdatedBy (optional) ---
	conn := config.GrpcClient()
	defer conn.Close()

	client := commonpb.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.GetUser(ctx, &commonpb.GetUserRequest{Id: int64(play.CreatedBy)})
	if err != nil {
		log.Println("Error calling GetUser:", err)
	}
	respupdate, err := client.GetUser(ctx, &commonpb.GetUserRequest{Id: int64(play.UpdatedBy)})
	if err != nil {
		log.Println("Error calling GetUser:", err)
	}

	return Resources.GameResource{
		ID:             play.ID,
		Status:         play.Status,
		BettingTime:    play.BettingTime,
		SettlementTime: play.SettlementTime,
		NoAutoExits:    play.NoAutoExits,
		GameDataID:     play.GameDataID,
		CreatedBy:      resp.Name,
		UpdatedBy:      respupdate.Name,
		CreatedAt:      play.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      play.UpdatedAt.Format("2006-01-02 15:04:05"),
		EnName:         getNameByLang(play.Translations, "en"),
		ZhName:         getNameByLang(play.Translations, "ch"),
		KmName:         getNameByLang(play.Translations, "kh"),
		GameTypeData:   groupedPlayTypeData, // grouped with GameTypeName filled
	}
}
func mapPlayLogToResource(play *models.Game) Resources.LogGameResource {
	// Map GameTypeData grouped by game_type_id
	playTypeDataMap := make(map[uint64][]Resources.GameTypeDataLogResource)

	for _, gtd := range play.GameTypeData {
		master := gtd.MasterGameType
		playTypeID := master.GameTypeID // <-- group by this

		image := ""
		if gtd.Image != "" {
			image = os.Getenv("MINIO_PUBLIC_ENDPOINT") + "/" + gtd.Image
		}

		resource := Resources.GameTypeDataLogResource{
			Odds:             gtd.Odds,
			NameKm:           gtd.NameKm,
			MasterGameTypeID: gtd.MasterGameTypeID,
			Image:            gtd.Image,
			ImageUrl:         image,
			KeyName:          master.KeyName,
			NameEn:           getPlayTypeNameByLang(master.GameTypeTranslation, "en"),
			NameZh:           getPlayTypeNameByLang(master.GameTypeTranslation, "ch"),
		}

		playTypeDataMap[playTypeID] = append(playTypeDataMap[playTypeID], resource)
	}

	// Convert map to slice for response
	groupedPlayTypeData := make([]Resources.GameTypeLogGroupResource, 0, len(playTypeDataMap))
	for playTypeID, items := range playTypeDataMap {
		playTypeName := ""

		// Find the actual GameTypeName from game.GameTypes
		for _, gt := range play.GameTypes {
			if gt.ID == playTypeID {
				playTypeName = gt.GameTypeName
				break
			}
		}

		sortPlayTypeDataLogResources(items)
		groupedPlayTypeData = append(groupedPlayTypeData, Resources.GameTypeLogGroupResource{
			GameTypeID:   playTypeID,
			GameTypeName: playTypeName, // now filled
			GameTypeData: items,
		})
	}

	sort.Slice(groupedPlayTypeData, func(i, j int) bool {
		return groupedPlayTypeData[i].GameTypeID < groupedPlayTypeData[j].GameTypeID
	})

	// --- GRPC calls for CreatedBy / UpdatedBy (optional) ---
	conn := config.GrpcClient()
	defer conn.Close()

	client := commonpb.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.GetUser(ctx, &commonpb.GetUserRequest{Id: int64(play.CreatedBy)})
	if err != nil {
		log.Println("Error calling GetUser:", err)
	}
	respupdate, err := client.GetUser(ctx, &commonpb.GetUserRequest{Id: int64(play.UpdatedBy)})
	if err != nil {
		log.Println("Error calling GetUser:", err)
	}

	return Resources.LogGameResource{
		ID:             play.ID,
		Status:         play.Status,
		BettingTime:    play.BettingTime,
		SettlementTime: play.SettlementTime,
		NoAutoExits:    play.NoAutoExits,
		CreatedBy:      resp.Name,
		UpdatedBy:      respupdate.Name,
		CreatedAt:      play.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      play.UpdatedAt.Format("2006-01-02 15:04:05"),
		EnName:         getNameByLang(play.Translations, "en"),
		ZhName:         getNameByLang(play.Translations, "ch"),
		KmName:         getNameByLang(play.Translations, "kh"),
		GameTypeData:   groupedPlayTypeData, // grouped with GameTypeName filled
	}
}
func sortPlayTypeDataLogResources(items []Resources.GameTypeDataLogResource) {
	sort.SliceStable(items, func(i, j int) bool {
		left := redLimitPlayOrder[strings.TrimSpace(strings.ToLower(items[i].KeyName))]
		right := redLimitPlayOrder[strings.TrimSpace(strings.ToLower(items[j].KeyName))]

		if left == 0 {
			left = 999
		}
		if right == 0 {
			right = 999
		}

		return left < right
	})
}

func sortPlayTypeDataModels(items []models.GameTypeData) {
	sort.SliceStable(items, func(i, j int) bool {
		left := 999
		right := 999

		if items[i].MasterGameType != nil {
			left = redLimitPlayOrder[strings.TrimSpace(
				strings.ToLower(items[i].MasterGameType.KeyName),
			)]
		}
		if items[j].MasterGameType != nil {
			right = redLimitPlayOrder[strings.TrimSpace(
				strings.ToLower(items[j].MasterGameType.KeyName),
			)]
		}

		if left == 0 {
			left = 999
		}
		if right == 0 {
			right = 999
		}

		return left < right
	})
}
func getNameByLang(translations []models.GameTranslation, lang string) string {
	for _, t := range translations {
		if t.LanguageCode == lang {
			return t.Content
		}
	}

	return ""
}
func getPlayTypeNameByLang(translations []models.GameTypeTranslation, lang string) string {
	for _, t := range translations {
		if t.LanguageCode == lang {
			return t.Name
		}
	}

	// fallback
	if len(translations) > 0 {
		return translations[0].Name
	}

	return ""
}

func mapPlayListResource(play *models.Game) Resources.GameListResource {

	return Resources.GameListResource{
		ID:             play.ID,
		BettingTime:    play.BettingTime,
		SettlementTime: play.SettlementTime,

		//Name:       getNameByLang(game.Translations, lang),
		EnName: getNameByLang(play.Translations, "en"),
		ZhName: getNameByLang(play.Translations, "ch"),
		KmName: getNameByLang(play.Translations, "kh"),
	}
}
func mapGamePlayTypeResource(playType *models.GameType) Resources.GamePlayTypeResource {
	return Resources.GamePlayTypeResource{
		ID:           playType.ID,
		GameTypeName: playType.GameTypeName,
	}
}
func mapPlayListFEResource(play *models.Game, lang string) Resources.GameListFEResource {

	return Resources.GameListFEResource{
		ID:             play.ID,
		BettingTime:    play.BettingTime,
		SettlementTime: play.SettlementTime,
		NoAutoExits:    play.NoAutoExits,
		Name:           getNameByLang(play.Translations, lang),
		Status:         play.Status,
		//EnName: getNameByLang(game.Translations, "en"),
		//ZhName: getNameByLang(game.Translations, "ch"),
		//KmName: getNameByLang(game.Translations, "kh"),
	}
}

// Display order only — does NOT affect PlayID values
var redLimitPlayOrder = map[string]int{
	"tiger":         1,
	"fish":          2,
	"shrimp":        3, // Gourd
	"crab":          4,
	"chicken":       5,
	"gold_coin":     6,
	"klouk":         7,
	"double_animal": 8,
	"triple_symbol": 9,  // Double Dices
	"triple_color":  10, // keep last (not in your list)
}

func sortPlayTypeDataResources(items []Resources.GameTypeDataResource) {
	sort.SliceStable(items, func(i, j int) bool {
		// left := redLimitPlayOrder[strings.TrimSpace(strings.ToLower(items[i].NameEn))]
		// right := redLimitPlayOrder[strings.TrimSpace(strings.ToLower(items[j].NameEn))]

		left := redLimitPlayOrder[strings.TrimSpace(strings.ToLower(items[i].KeyName))]
		right := redLimitPlayOrder[strings.TrimSpace(strings.ToLower(items[j].KeyName))]

		// fallback if NameEn isn't a key — use MasterGameTypeID or keep stable
		if left == 0 {
			left = 999
		}
		if right == 0 {
			right = 999
		}

		return left < right
	})
}

// GetGames godoc
// @Summary      List Games
// @Description  Get all games
// @Tags         Games
// @Produce      json
// @Success      200 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /games [get]
func (gc *PlayController) GetGames(c *fiber.Ctx) error {
	// Get all games from service

	//lang := c.Get("Accept-Language", "en")
	plays, err := gc.Service.All()
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	// Map games to resources
	resp := make([]Resources.GameResource, 0, len(plays))
	for i := range plays {
		resp = append(resp, mapPlayToResource(&plays[i]))
	}

	// Return plain JSON without pagination
	return ApiResponse.SuccessWithPagination(c, 200, len(plays), resp, "Games retrieved successfully")
}

// CreateGame godoc
// @Summary      Create a Game
// @Description  Create a new game
// @Tags         Games
// @Accept       json
// @Produce      json
// @Param        payload body requests.GameCreateRequest true "Game Data"
// @Success      201 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /games [post]
func (gc *PlayController) CreatePlay(c *fiber.Ctx) error {
	req := new(requests.GameCreateRequest)

	user, _ := c.Locals("user").(jwt.MapClaims)

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
	play, err := gc.Service.CreatePlay(*req, uint(user["id"].(float64)))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, "fail")
	}
	// newLog := map[string]any{
	// 	"status":     gamehall.Status,
	// 	"sort_order": gamehall.SortOrder,
	// 	"en_name":    getNameByLang(gamehall.Translations, "en"),
	// 	"ch_name":    getNameByLang(gamehall.Translations, "ch"),
	// 	"kh_name":    getNameByLang(gamehall.Translations, "kh"),
	// }
	// gc.Service.AuditAction(
	// 	c,
	// 	"create",
	// 	"game",
	// 	user,
	// 	nil,
	// 	game,
	// )
	resp := mapPlayLogToResource(play)
	gc.Service.AuditAction(
		c,
		"create",
		"game",
		user,
		nil,
		resp,
	)
	return ApiResponse.Success(c, fiber.StatusCreated, resp, "Game created successfully")
}

// UpdateGame godoc
// @Summary      Update a Game
// @Description  Update an existing game
// @Tags         Games
// @Accept       json
// @Produce      json
// @Param        id path int true "Game ID"
// @Param        payload body requests.GameUpdateRequest true "Game Data"
// @Success      200 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /games/{id} [put]
func (gc *PlayController) UpdatePlay(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id < 0 {
		return ApiResponse.ValidationError(c, fiber.StatusBadRequest, "Invalid ID", err.Error())
	}
	if id < 0 {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "Invalid ID")
	}
	req := new(requests.GameUpdateRequest)

	if err := c.BodyParser(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("bad request", err.Error())
		return ApiResponse.ValidationError(c, fiber.StatusBadRequest, "unable to parse request body", fullMessages[0])
	}

	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("validation error", fullMessages[0])
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
	playold, err := gc.Service.FindByID(uint64(id))
	play, err := gc.Service.UpdatePlay(uint64(id), *req, uint(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ApiResponse.Error(c, fiber.StatusNotFound, "record not found")
		}

		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	resp := mapPlayLogToResource(play)
	oldresp := mapPlayLogToResource(playold)

	//oldData, newData := utils.StructChanges(oldresp, resp)

	gc.Service.AuditAction(
		c,
		"update",
		"game",
		user,
		oldresp,
		resp,
	)
	return ApiResponse.Success(c, fiber.StatusOK, resp, "Game updated successfully")
}

// DeleteGame godoc
// @Summary      Delete a Game
// @Description  Delete an existing game
// @Tags         Games
// @Accept       json
// @Param        id path int true "Game ID"
// @Success      200 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /games/{id} [delete]
func (gc *PlayController) DeletePlay(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "Invalid ID")
	}
	if id < 0 {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "Invalid ID")
	}
	playold, err := gc.Service.Repo.FindByID(uint64(id))
	err = gc.Service.DeletePlay(uint64(id))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}
	oldresp := mapPlayLogToResource(playold)
	gc.Service.AuditAction(
		c,
		"delete",
		"game",
		c.Locals("user").(jwt.MapClaims),
		oldresp,
		nil,
	)

	return ApiResponse.Success(c, 200, nil, "Game deleted successfully")
}

// GetGames godoc
// @Summary      List Games to use in game table
// @Description  Get all games
// @Tags         Games
// @Produce      json
// @Success      200 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /games/list [get]
func (gc *PlayController) GetPlaysList(c *fiber.Ctx) error {
	// Get all games from service

	//lang := c.Get("Accept-Language", "en")
	plays, err := gc.Service.AllEnable()
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	// Map games to resources
	resp := make([]Resources.GameListResource, 0, len(plays))
	for i := range plays {
		resp = append(resp, mapPlayListResource(&plays[i]))
	}

	// Return plain JSON without pagination
	return ApiResponse.SuccessWithPagination(c, 200, len(plays), resp, "Games retrieved successfully")
}

func (gc *PlayController) GetPlayTypeByPlay(c *fiber.Ctx) error {
	id, err := c.ParamsInt("game_id")
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "Invalid ID")
	}

	if id < 0 {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "Invalid ID")
	}

	playtypes, err := gc.Service.GetPlayTypesByPlayID(uint64(id))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusNotFound, err.Error())
	}

	// Map games to resources
	resp := make([]Resources.GamePlayTypeResource, 0, len(playtypes))
	for i := range playtypes {
		resp = append(resp, mapGamePlayTypeResource(&playtypes[i]))
	}

	return ApiResponse.Success(c, 200, resp, "Game types fetched successfully")
}

// GetGames godoc
// @Summary      List Games to use in game table
// @Description  Get all games
// @Tags         Games
// @Produce      json
// @Success      200 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /public/api/game-list [get]
func (gc *PlayController) GetPlaysListForFE(c *fiber.Ctx) error {
	// Get all games from service

	lang := c.Get("Accept-Language", "en")
	plays, err := gc.Service.AllEnable()
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	// Map games to resources
	resp := make([]Resources.GameListFEResource, 0, len(plays))
	for i := range plays {
		resp = append(resp, mapPlayListFEResource(&plays[i], lang))
	}

	// Return plain JSON without pagination
	return ApiResponse.Success(c, 200, resp, "Games retrieved successfully")
}
