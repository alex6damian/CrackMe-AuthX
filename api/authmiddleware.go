package api

import (
	"strings"

	"github.com/alex6damian/CrackMe-AuthX/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT token and sets user context
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	var token string

	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			token = parts[1]
		}
	}

	// fallback la cookie
	if token == "" {
		token = c.Cookies("token")
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Missing token",
		})
	}

	if utils.IsTokenBlacklisted(token) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Token has been invalidated",
		})
	}

	claims, err := utils.ValidateToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid or expired token",
		})
	}

	c.Locals("userID", claims.UserID)
	c.Locals("userEmail", claims.Email)
	c.Locals("userRole", claims.Role)
	return c.Next()
}
