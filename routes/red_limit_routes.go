package routes

import (
	"config-service/config"
	"config-service/controllers"
	"config-service/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func RedLimitRoutes(router fiber.Router, redisConn *redis.Client) {
	gtc := controllers.NewPlayRedLimitController(config.DB, redisConn)
	gt := router.Group("/red-limit", middlewares.Protected())
	gt.Post("/", gtc.CreateRedLimit)
	gt.Get("/", gtc.GetRedLimits)
	gt.Get("/:id", gtc.GetRedLimit)
	gt.Put("/:id", gtc.UpdateRedLimit)
	gt.Delete("/:id", gtc.DeleteRedLimit)

}
