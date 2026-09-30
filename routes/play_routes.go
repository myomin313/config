package routes

import (
	"config-service/config"
	"config-service/controllers"
	"config-service/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func PlayRoutes(router fiber.Router, redisConn *redis.Client) {
	gc := controllers.NewPlayController(config.DB, redisConn)
	play := router.Group("/games", middlewares.Protected())
	play.Get("/", gc.GetGames)
	play.Post("/", gc.CreatePlay)
	play.Put("/:id", gc.UpdatePlay)
	play.Delete("/:id", gc.DeletePlay)
	play.Get("/external/list", gc.GetPlaysList)
	router.Get("/get-game-types/:game_id", gc.GetPlayTypeByPlay)

	router.Get("/game-list", gc.GetPlaysListForFE)

}
