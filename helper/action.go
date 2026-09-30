package helper

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mssola/user_agent"
)

func PublishAudit(
	auditService *AuditService,
	c *fiber.Ctx,
	action string,
	model string,
	user jwt.MapClaims,
	before any,
	after any,
) {

	if auditService == nil {
		log.Println("Audit service not initialized")
		return
	}

	if auditService.Channel == nil {
		log.Println("RabbitMQ channel not initialized")
		return
	}

	ua := user_agent.New(c.Get("User-Agent"))

	device := "Desktop"
	if ua.Mobile() {
		device = "Mobile"
	}

	audit := map[string]any{
		"action":     action,
		"model":      model,
		"ip_address": GetClientIP(c),
		"device":     device,
		"created_at": time.Now(),
		"role":       user["role"],
		"user": map[string]any{
			"email": user["email"],
		},
	}

	if before != nil {
		audit["old_data"] = before
	}

	if after != nil {
		audit["new_data"] = after
	}

	auditService.Publish(model, audit)
}
func GetClientIP(c *fiber.Ctx) string {

	// Check proxy headers first
	if ip := c.Get("X-Forwarded-For"); ip != "" {
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	if ip := c.Get("X-Real-IP"); ip != "" {
		return ip
	}

	return c.IP()
}

// func PublishAudit(
// 	auditService *AuditService,
// 	c *fiber.Ctx,
// 	action string,
// 	model string,
// 	user jwt.MapClaims,
// 	before any,
// 	after any,
// ) {
// 	if auditService == nil {
// 		log.Println("❌ Audit service not initialized")
// 		return
// 	}

// 	if auditService.Redis == nil {
// 		log.Println("❌ Redis client not initialized")
// 		return
// 	}

// 	ua := user_agent.New(c.Get("User-Agent"))

// 	device := "Desktop"
// 	if ua.Mobile() {
// 		device = "Mobile"
// 	}

// 	audit := map[string]any{
// 		"action":     action,
// 		"model":      model,
// 		"ip_address": c.IP(),
// 		"device":     device,
// 		"created_at": time.Now(),
// 		"role":       user["role"],
// 		"user": map[string]any{
// 			"email": user["email"],
// 		},
// 	}

// 	if before != nil {
// 		audit["old_data"] = before
// 	}

// 	if after != nil {
// 		audit["new_data"] = after
// 	}

// 	auditService.Publish(model, audit, time.Minute)
// }
