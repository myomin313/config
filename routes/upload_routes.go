package routes

import (
	"config-service/controllers"
	"config-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

func UploadRoutes(router fiber.Router) {
	router.Post("/external/uploadfile", middlewares.Protected(), controllers.UploadToMinio)
	router.Post("/external/uploadfiles-rule", middlewares.Protected(), controllers.UploadToMinioplayrule)
}
