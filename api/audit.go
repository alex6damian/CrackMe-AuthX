package api

import (
	"time"

	"github.com/alex6damian/CrackMe-AuthX/internal/models"
	"github.com/alex6damian/CrackMe-AuthX/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GET /api/audit  (manager only)
func ListAuditLogs(c *fiber.Ctx) error {
	logs, err := repository.GetAllAuditLogs()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching audit logs"})
	}
	return c.JSON(fiber.Map{"success": true, "data": logs})
}

// logAction inserts an audit log entry — called internally by handlers
func logAction(c *fiber.Ctx, userID uint, action, resource, resourceID string) {
	entry := models.AuditLog{
		ID:         uuid.NewString(),
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Timestamp:  time.Now(),
		IPAddress:  c.IP(),
	}
	_ = repository.CreateAuditLog(&entry)
}
