package controllers

import (
	"config-service/ApiResponse"
	"config-service/Repositories"
	Resources "config-service/Resources"
	"config-service/Services"
	"config-service/models"

	//"config-service/requests"
	// "config-service/utils"
	// "log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PlayTypeController struct {
	Service *Services.PlayTypeService
}

func NewPlayTypeController(db *gorm.DB) *PlayTypeController {
	return &PlayTypeController{
		Service: &Services.PlayTypeService{
			Repo: &Repositories.PlayTypeRepository{
				DB: db,
			},
		},
	}
}

func mapPlayTypeToResource(g *models.GameType) Resources.GameTypeResource {

	return Resources.GameTypeResource{
		ID:           g.ID,
		GameTypeName: g.GameTypeName,
	}
}

// GetGames godoc
// @Summary      List GameTypes
// @Description  Get all GameTypes
// @Tags         GameTypes
// @Produce      json
// @Success      200 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /game-types [get]
func (gc *PlayTypeController) GetPlayTypes(c *fiber.Ctx) error {

	play_data_id, err := c.ParamsInt("game_data_id")
	if err != nil || play_data_id < 0 {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid_id")
	}

	GameTypes, err := gc.Service.All(uint(play_data_id))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	// Map games to resources
	resp := make([]Resources.GameTypeResource, 0, len(GameTypes))
	for i := range GameTypes {
		resp = append(resp, mapPlayTypeToResource(&GameTypes[i]))
	}

	// Return plain JSON without pagination
	return ApiResponse.SuccessWithPagination(c, 200, len(GameTypes), resp, "GameTypes retrieved successfully")
}
