package routes

import (
	"config-service/config"
	"config-service/controllers"
	"config-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

func ReportRoutes(router fiber.Router) {
	gc := controllers.NewReportController(config.DB)
	play := router.Group("/reports", middlewares.Protected())
	play.Get("/game-result", gc.GetPlayResult)
	play.Put("/round/:id/remark", gc.UpdateRoundRemark)
	play.Put("/round/:id/result", gc.UpdateRoundResult)
}
