package middlewares

import (
	"config-service/ApiResponse"
	"config-service/utils"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("APP_KEY"))

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")

		if tokenString == "" {
			return ApiResponse.ValidationError(c, fiber.StatusUnauthorized, "Missing token")
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		decryptedToken, err := utils.Decrypt(tokenString)
		if err != nil {
			return ApiResponse.ValidationError(c, fiber.StatusUnauthorized, "invalid encrypted token")
		}

		// token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// 	return jwtSecret, nil
		// })
		token, err := jwt.Parse(
			decryptedToken,
			func(t *jwt.Token) (interface{}, error) {

				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fiber.NewError(
						fiber.StatusUnauthorized,
						"Invalid signing method",
					)
				}

				return jwtSecret, nil
			},
		)

		if err != nil || !token.Valid {
			return ApiResponse.ValidationError(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			return ApiResponse.ValidationError(
				c,
				fiber.StatusUnauthorized,
				"Invalid token claims",
			)
		}
		log.Print("admin middleware====", claims)
		// 🔓 Decrypt only user ID
		if encryptedID, ok := claims["id"].(string); ok {
			log.Print("encryptedID string", encryptedID)
			decryptedID, err := utils.Decrypt(encryptedID)
			if err != nil {
				return ApiResponse.ValidationError(
					c,
					fiber.StatusUnauthorized,
					"Failed to decrypt user id",
				)
			}

			// convert decrypted string -> uint64
			idUint64, err := strconv.ParseUint(decryptedID, 10, 64)
			if err != nil {
				return ApiResponse.ValidationError(
					c,
					fiber.StatusUnauthorized,
					"Invalid user id format",
				)
			}

			// IMPORTANT: controller expects float64
			claims["id"] = float64(idUint64)
		}

		// Store user in context
		c.Locals("user", claims)

		//c.Locals("user", token.Claims)
		return c.Next()
	}
}
