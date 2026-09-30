package routes

import (
	"config-service/config"
	"config-service/controllers"
	"config-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

func MaintenanceRoutes(router fiber.Router) {
	mc := controllers.NewMaintenanceController(config.DB)
	m := router.Group("/maintenance", middlewares.Protected())
	m.Get("/", mc.GetMaintenance)
	m.Put("/", mc.UpdateMaintenance)

	router.Get("/application/setting", mc.GetGetMaintenanceForFE)
}
