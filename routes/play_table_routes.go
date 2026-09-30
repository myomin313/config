package routes

import (
	"config-service/config"
	"config-service/controllers"
	"config-service/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func PlayTableRoutes(router fiber.Router, redisConn *redis.Client) {
	gtc := controllers.NewPlayTableController(config.DB, redisConn)
	gt := router.Group("/gametables", middlewares.Protected())
	gt.Post("/", gtc.CreatePlaytable)
	gt.Post("/status", gtc.MassUpdatePlayTableStatus)
	gt.Get("/", gtc.GetPlayTables)
	gt.Get("/:id", gtc.GetPlayTable)
	gt.Put("/:id", gtc.UpdatePlayTable)
	gt.Delete("/:id", gtc.DeletePlayTable)

	router.Get("/gametable-list", gtc.GetPlayTableForFE)

}
