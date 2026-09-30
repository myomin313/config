package routes

import (
	"config-service/config"
	"config-service/controllers"

	"github.com/gofiber/fiber/v2"
)

func RuleRoutes(router fiber.Router) {
	rc := controllers.NewRuleController(config.DB)
	rule := router.Group("/rules")
	rule.Get("/get-rule", rc.GetRules)
}
