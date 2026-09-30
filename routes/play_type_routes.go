package routes

import (
	"config-service/config"
	"config-service/controllers"

	"github.com/gofiber/fiber/v2"
)

func PlayTypeRoutes(router fiber.Router) {
	gc := controllers.NewPlayTypeController(config.DB)
	play := router.Group("/game-types")
	play.Get("/:game_data_id", gc.GetPlayTypes)

}
