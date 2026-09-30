package routes

import (
	"config-service/config"
	"config-service/controllers"
	"config-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

func AnnouncementRoutes(router fiber.Router) {
	ac := controllers.NewAnnouncementController(config.DB)
	announcement := router.Group("/announcements", middlewares.Protected())
	announcement.Get("/", ac.GetAnnouncementsList)
	announcement.Post("/", ac.CreateAnnouncement)
	announcement.Delete("/multiple/delete", ac.DeleteMultipleAnnouncement)
	announcement.Put("/:id", ac.UpdateAnnouncement)
	announcement.Delete("/:id", ac.DeleteAnnouncement)
	announcement.Post("/status", ac.MassUpdateAnnouncementStatus)
	announcement.Post("/:id/detail", ac.CreateAnnouncementDetail)
	announcement.Put("/detail/update/:id", ac.UpdateAnnouncementDetail)
	announcement.Delete("/detail/delete/:id", ac.DeleteAnnouncementDetail)
	announcement.Get("/:id", ac.GetAnnouncementsDetail)
	//game:= router.Group("/announcements-fe")
	router.Get("/announcement-list", ac.GetAnnouncementsListForFE)

	// gamehall.Get("/list", ac.GetGameHallsList)
}
