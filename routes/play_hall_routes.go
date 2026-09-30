package routes

import (
	"config-service/config"
	"config-service/controllers"
	"config-service/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func PlayHallRoutes(router fiber.Router, redisConn *redis.Client) {
	gc := controllers.NewPlayHallController(config.DB, redisConn)
	playhall := router.Group("/gamehalls", middlewares.Protected())
	playhall.Get("/", gc.GetPlayHalls)
	playhall.Post("/", gc.CreatePlayHall)
	playhall.Put("/:id", gc.UpdatePlayHall)
	playhall.Delete("/:id", gc.DeletePlayHall)
	playhall.Get("/external/list", gc.GetPlayHallsList)

	router.Get("/gamehall-list", gc.GetPlayHallsListForFE)
}
