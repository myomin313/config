package routes

import (
	"config-service/config"
	"config-service/controllers"

	"github.com/gofiber/fiber/v2"
)

func BannerRoutes(router fiber.Router) {
	controller := controllers.NewBannerController(config.DB)
	b := router.Group("/banners")
	b.Post("/", controller.CreateBanner)
	b.Post("/status", controller.MassUpdateBannerStatus)
	b.Post("/single-status", controller.UpdateBannerStatus)
	b.Get("/", controller.GetAllBanners)
	b.Get("/:id", controller.GetBannerByID)
	b.Put("/:id", controller.UpdateBanner)
	b.Delete("/:id", controller.DeleteBanner)
	b.Post("/mass-delete", controller.MassBannerDelete)

	router.Post("/banner-list", controller.GetActiveBanners)
}
