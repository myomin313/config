package routes

import (
	"config-service/config"
	"config-service/controllers"

	"github.com/gofiber/fiber/v2"
)

func PlayDataRoutes(router fiber.Router) {
	gc := controllers.NewPlayDataController(config.DB)
	play := router.Group("/game-data")
	play.Get("/", gc.GetPlayData)
}
