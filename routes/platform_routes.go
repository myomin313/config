package routes

import (
	"config-service/config"
	"config-service/controllers"
	"config-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

func PlatformRoutes(router fiber.Router) {
	pc := controllers.NewPlatformController(config.DB)
	platform := router.Group("/platforms", middlewares.Protected())
	platform.Post("/", pc.CreatePlatform)
	platform.Get("/", pc.GetPlatforms)
	platform.Get("/:id", pc.GetPlatform)
	platform.Put("/:id", pc.UpdatePlatform)
	platform.Delete("/:id", pc.DeletePlatform)
	platform.Get("/merchant/currencies", pc.GetMerchantCurrency)
	platform.Get("/:id/merchants", pc.GetMerchantList)
	platform.Post("/:id/merchants", pc.CreateMerchantCurrency)
	platform.Put("/merchant-currencies/:id", pc.UpdateMerchantCurrency)
	platform.Get("/all/platforms", pc.GetAllPlatforms)
}
