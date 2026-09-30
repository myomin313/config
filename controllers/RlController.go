package controllers

import (
	"strings"

	"config-service/ApiResponse"
	Resources "config-service/Resources"
	"config-service/Services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ------------------------ Controller Struct ------------------------

type RuleController struct {
	Service *Services.RuleService
}

func NewRuleController(db *gorm.DB) *RuleController {
	return &RuleController{
		Service: &Services.RuleService{},
	}
}

// ------------------------ Handlers ------------------------

// GetRules godoc
// @Summary      Get Rule Image URL
// @Description  Get rule image URL based on game parameters
// @Tags         Rules
// @Produce      json
// @Param        game_id query string true "Game ID"
// @Param        game_type query string true "Game Type"
// @Param        agent query string true "Agent"
// @Param        language_code query string true "Language Code"
// @Success      200 {object} Resources.RuleResource
// @Failure      400 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /rules/get-rule [get]
func (rc *RuleController) GetRules(c *fiber.Ctx) error {
	// Get query parameters
	playID := strings.TrimSpace(c.Query("game_id"))
	playType := strings.TrimSpace(c.Query("game_type"))
	agent := strings.TrimSpace(c.Query("agent"))
	languageCode := c.Get("Accept-Language", "en")
	// Validate required parameters
	if playID == "" {
		return ApiResponse.Error(c, fiber.StatusUnprocessableEntity, "Missing required parameters: game_id")
	}
	if playType == "" {
		return ApiResponse.Error(c, fiber.StatusUnprocessableEntity, "Missing required parameters: game_type")
	}
	// if agent == ""{
	// 	return ApiResponse.Error(c, fiber.StatusUnprocessableEntity, "Missing required parameters: agent")
	// }
	if languageCode == "" {
		return ApiResponse.Error(c, fiber.StatusUnprocessableEntity, "Missing required headers: Accept-Language")
	}

	// Call service to construct URL
	url, err := rc.Service.GetRuleImageURL(playID, playType, agent, languageCode)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	// Map to response resource
	resp := Resources.RuleResource{
		URL: url,
	}

	//return ApiResponse.Success(c, 200, resp, "Rule URL retrieved successfully")
	return ApiResponse.JsonSuccess(c, 200, resp, "Rule URL retrieved successfully")

}
