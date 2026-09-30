package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx) error {
	return c.SendString("welcome to kla-klouk-game")
}

func Health(c *fiber.Ctx) error {
	return c.SendStatus(200) // explicitly
}

func Test(c *fiber.Ctx) error {
	response := fiber.Map{
		"platform": fiber.Map{
			"id":     1,
			"name":   "Kla Klouk Platform",
			"email":  "admin@klaklouk.com",
			"status": "active",
			"api_key": fiber.Map{
				"id":         10,
				"key":        "pk_live_123456",
				"status":     "enabled",
				"created_at": "2023-01-01T00:00:00Z",
				"updated_at": "2023-01-01T00:00:00Z"},
		},
	}

	return c.JSON(response)
}
