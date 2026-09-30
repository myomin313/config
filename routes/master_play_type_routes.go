package routes

import (
	"config-service/config"
	"config-service/controllers"

	"github.com/gofiber/fiber/v2"
)

func MasterPlayTypeRoutes(router fiber.Router) {
	gc := controllers.NewMasterPlayTypeController(config.DB)
	play := router.Group("/master-game-types")
	play.Get("/", gc.GetMasterPlayTypes)

}
