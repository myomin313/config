package controllers

import (
	"config-service/ApiResponse"
	"config-service/Repositories"
	"config-service/Services"
	"config-service/config"
	"config-service/helper"
	"config-service/models"
	"config-service/requests"
	"config-service/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type MaintenanceController struct {
	Service *Services.MaintenanceService
}

func NewMaintenanceController(db *gorm.DB) *MaintenanceController {
	auditService := &helper.AuditService{
		Channel: config.RabbitChannel,
		Queue:   "audit_logs",
	}

	return &MaintenanceController{
		Service: &Services.MaintenanceService{
			Repo: &Repositories.MaintenanceRepository{
				DB: db,
			},
			Audit: auditService,
		},
	}
}

func mapMaintenanceToResponse(p *models.Maintenance) map[string]interface{} {
	resp := map[string]interface{}{
		"id":                      p.ID,
		"application_maintenance": p.ApplicationMaintenance,
		"game_maintenance":        p.GameMaintenance,
		"created_at":              p.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at":              p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return resp
}

func mapMaintenanceToResponseFE(p *models.Maintenance) map[string]interface{} {
	resp := map[string]interface{}{
		"applicationMaintenanceSwitch": p.ApplicationMaintenance,
		"gameMaintenanceSwitch":        p.GameMaintenance,
	}
	return resp
}

// GetMaintenance godoc
// @Summary      Get maintenance flags
// @Description  Retrieve current maintenance settings
// @Tags         Maintenance
// @Produce      json
// @Success      200 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /maintenance [get]
func (pc *MaintenanceController) GetMaintenance(c *fiber.Ctx) error {
	m, err := pc.Service.Get()
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	resp := mapMaintenanceToResponse(m)
	return ApiResponse.Success(c, 200, resp, "maintenance_retrieved_successfully")
}

// GetGetMaintenanceForFE godoc
// @Summary      Get maintenance flags
// @Description  Retrieve current maintenance settings
// @Tags         Maintenance
// @Produce      json
// @Success      200 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /application/setting [get]
func (pc *MaintenanceController) GetGetMaintenanceForFE(c *fiber.Ctx) error {
	m, err := pc.Service.Get()
	if err != nil {
		return ApiResponse.JsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	resp := mapMaintenanceToResponseFE(m)
	return ApiResponse.JsonSuccess(c, 200, resp, "success")
}

// UpdateMaintenance godoc
// @Summary      Update maintenance flags
// @Description  Update maintenance settings
// @Tags         Maintenance
// @Accept       json
// @Produce      json
// @Param        payload body requests.MaintenanceUpdateRequest true "Maintenance Data"
// @Success      200 {object} ApiResponse.Response
// @Failure      400 {object} ApiResponse.Response
// @Failure      401 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /maintenance [put]
func (pc *MaintenanceController) UpdateMaintenance(c *fiber.Ctx) error {
	req := new(requests.MaintenanceUpdateRequest)

	// Parse JSON body
	if err := c.BodyParser(req); err != nil {
		// Fallback to form-data parsing
		appMaint := c.FormValue("application_maintenance")
		playMaint := c.FormValue("game_maintenance")

		if appMaint != "" {
			val := appMaint == "true"
			req.ApplicationMaintenance = &val
		}
		if playMaint != "" {
			val := playMaint == "true"
			req.GameMaintenance = &val
		}
	}

	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation_error", fullMessages[0])
	}

	user, ok := c.Locals("user").(jwt.MapClaims)
	if !ok || user == nil {
		return ApiResponse.Error(c, fiber.StatusUnauthorized, "user_not_authenticated")
	}
	userIDf, ok := user["id"].(float64)
	if !ok {
		return ApiResponse.Error(c, fiber.StatusUnauthorized, "invalid_user_id_in_token")
	}
	oldm, err := pc.Service.Get()
	oldresp := mapMaintenanceToResponse(oldm)
	m, err := pc.Service.Update(req, uint(userIDf))
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	resp := mapMaintenanceToResponse(m)
	pc.Service.AuditAction(
		c,
		"update",
		"Maintanance",
		user,
		oldresp,
		resp,
	)
	return ApiResponse.Success(c, 200, resp, "maintenance_updated_successfully")
}
