package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func Setup(app *fiber.App, redisConn *redis.Client) {
	api := app.Group("/api")
	Welcome(app)
	Testing(api)
	// PlatformRoutes(api)
	PlayRoutes(api, redisConn)
	MasterPlayTypeRoutes(api)
	PlayTableRoutes(api, redisConn)
	UploadRoutes(api)
	PlayTypeRoutes(api)
	RedLimitRoutes(api, redisConn)
	PlayHallRoutes(api, redisConn)
	AnnouncementRoutes(api)
	BannerRoutes(api)
	PlatformRoutes(api)
	PlayDataRoutes(api)
	RuleRoutes(api)
	PlayResultRoutes(api)
	ReportRoutes(api)
	MaintenanceRoutes(api)
}
