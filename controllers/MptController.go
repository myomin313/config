package controllers

import (
	"config-service/ApiResponse"
	"config-service/Repositories"
	Resources "config-service/Resources"
	"config-service/Services"
	"config-service/models"
	"config-service/utils"

	//"config-service/requests"
	// "config-service/utils"
	"log"
	"sort"

	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Title struct {
	ID      string `json:"id"`
	Odds    string `json:"odds"`
	Image   string `json:"image"`
	KeyName string `json:"key_name"`
	NameEn  string `json:"name_en"`
	NameKm  string `json:"name_km"`
	NameZh  string `json:"name_zh"`
}

type MasterPlayTypeController struct {
	Service *Services.MasterPlayTypeService
}

func NewMasterPlayTypeController(db *gorm.DB) *MasterPlayTypeController {
	return &MasterPlayTypeController{
		Service: &Services.MasterPlayTypeService{
			Repo: &Repositories.MasterPlayTypeRepository{
				DB: db,
			},
		},
	}
}

func mapMasterPlayTypeToResource(g *models.MasterGameType) Resources.MasterGameTypeResource {

	return Resources.MasterGameTypeResource{
		ID:      g.ID,
		Odds:    g.Odds,
		Image:   g.Image,
		KeyName: g.KeyName,
		NameEn:  getPlayTypeNameByLang(g.GameTypeTranslation, "en"),
		NameKm:  getPlayTypeNameByLang(g.GameTypeTranslation, "kh"),
		NameZh:  getPlayTypeNameByLang(g.GameTypeTranslation, "ch"),
	}
}

// Get all master game types
// @Summary Get all master game types
// @Description Get all master game types
// @Tags MasterGameTypes
// @Produce json
// @Success 200 {object} ApiResponse.Response
// @Failure 500 {object} ApiResponse.Response
// @Router /master-game-types [get]
func (gc *MasterPlayTypeController) GetMasterPlayTypes(c *fiber.Ctx) error {

	// 🔹 Parse query param: ?game_type_id=1,2,3
	playTypeIDsStr := c.Query("game_type_id")
	var playTypeIDs []uint64

	if playTypeIDsStr != "" {
		ids := strings.Split(playTypeIDsStr, ",")
		for _, idStr := range ids {
			idStr = strings.TrimSpace(idStr)
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				return ApiResponse.Error(c, fiber.StatusBadRequest, "Invalid game_type_id: "+idStr)
			}
			playTypeIDs = append(playTypeIDs, id)
		}
	}

	log.Println("game type filter:", playTypeIDs)

	// 🔹 Fetch from service
	masterPlayTypes, err := gc.Service.All(playTypeIDs...)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	// 🔹 Group by GameTypeID
	groupMap := make(map[uint64]*Resources.GroupedMasterGameType)

	for i := range masterPlayTypes {

		resource := mapMasterPlayTypeToResource(&masterPlayTypes[i])
		log.Println("masterGameTypes", masterPlayTypes[i].GameTypeID)
		playTypeID := masterPlayTypes[i].GameTypeID
		playTypeName := masterPlayTypes[i].GameType.GameTypeName // MUST preload GameType

		if _, exists := groupMap[playTypeID]; !exists {
			groupMap[playTypeID] = &Resources.GroupedMasterGameType{
				GameTypeID:   playTypeID,
				GameTypeName: playTypeName,
				Items:        []Resources.MasterGameTypeResource{},
			}
		}

		groupMap[playTypeID].Items = append(groupMap[playTypeID].Items, resource)
	}

	// 🔹 Convert map → slice
	// groupedResp := make([]Resources.GroupedMasterGameType, 0, len(groupMap))
	// for _, group := range groupMap {

	// 	// Optional: sort items by ID
	// 	sort.Slice(group.Items, func(i, j int) bool {
	// 		return group.Items[i].ID < group.Items[j].ID
	// 	})

	// 	groupedResp = append(groupedResp, *group)
	// }

	// // Optional: sort groups by GameTypeID
	// sort.Slice(groupedResp, func(i, j int) bool {
	// 	return groupedResp[i].GameTypeID < groupedResp[j].GameTypeID
	// })
	// 🔹 Convert map → slice
	groupedResp := make([]Resources.GroupedMasterGameType, 0, len(groupMap))
	for _, group := range groupMap {

		// Sort items by play order (does NOT change IDs)
		sort.SliceStable(group.Items, func(i, j int) bool {
			return utils.PlayOrderIndex(group.Items[i].KeyName) <
				utils.PlayOrderIndex(group.Items[j].KeyName)
		})

		groupedResp = append(groupedResp, *group)
	}

	// Sort groups by GameTypeID (unchanged)
	sort.Slice(groupedResp, func(i, j int) bool {
		return groupedResp[i].GameTypeID < groupedResp[j].GameTypeID
	})

	// Define column title metadata
	title := Title{
		ID:      "ID",
		Odds:    "Odds",
		Image:   "Image",
		KeyName: "Key name",
		NameEn:  "Game name (English)",
		NameKm:  "Game name(កម្ពុជា)",
		NameZh:  "Game name(中文简体)",
	}

	return ApiResponse.Success(c, fiber.StatusOK, fiber.Map{
		"data":  groupedResp,
		"title": title,
	}, "masterGameTypes retrieved successfully")
}
