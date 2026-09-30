package controllers

import (
	//"time"

	"config-service/ApiResponse"
	"config-service/Repositories"
	Resources "config-service/Resources"
	"config-service/Services"
	"config-service/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ------------------------ Controller Struct ------------------------

type PlayDataController struct {
	Service *Services.PlayDataService
}

func NewPlayDataController(db *gorm.DB) *PlayDataController {
	return &PlayDataController{
		Service: &Services.PlayDataService{
			Repo: &Repositories.PlayDataRepository{
				DB: db,
			},
		},
	}
}

// ------------------------ Helper ------------------------

func mapPlayDataToResponse(p *models.GameData) Resources.GameDataResource {
	return Resources.GameDataResource{
		ID:     p.ID,
		NameCh: p.NameCh,
		NameEn: p.NameEn,
	}
}

// GetGameData godoc
// @Summary      List Games Data
// @Description  Get all Games Data
// @Tags         gameData
// @Produce      json
// @Success      200 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /game-data [get]
func (pc *PlayDataController) GetPlayData(c *fiber.Ctx) error {
	platforms, err := pc.Service.All()
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	resp := make([]Resources.GameDataResource, 0, len(platforms))
	for i := range platforms {
		resp = append(resp, mapPlayDataToResponse(&platforms[i]))
	}

	return ApiResponse.Success(c, 200, resp, "game data retrieved successfully")
}
