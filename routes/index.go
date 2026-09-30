package routes

import (
	"config-service/controllers"

	"github.com/gofiber/fiber/v2"
)

func Welcome(app *fiber.App) {
	app.Get("/", controllers.Index)
	app.Get("/health", controllers.Health)

}

// Testing godoc
// @Summary Get json response summary
// @Description Get json response Description
// @Tags Testing
// @Accept json
// @Produce json
// @Failure 404 {object} map[string]string
// @Router /testing [get]
func Testing(router fiber.Router) {
	router.Get("/testing", controllers.Test)
}

// GameResultRoutes godoc
// @Summary Get game result report
// @Description Get game result report with table columns as fields
// @Tags Game Result
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/game-result-report [get]
func PlayResultRoutes(router fiber.Router) {
	router.Get("/game-result-report", controllers.GetPlayResultReport)
}
