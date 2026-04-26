package api

import (
	"time"

	"github.com/alex6damian/CrackMe-AuthX/internal/models"
	"github.com/alex6damian/CrackMe-AuthX/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateTicketRequest struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Severity    models.Severity `json:"severity"`
}

type UpdateTicketRequest struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Severity    models.Severity `json:"severity"`
}

type ChangeStatusRequest struct {
	Status models.TicketStatus `json:"status"`
}

// POST /api/tickets
func CreateTicket(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req CreateTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Error parsing"})
	}

	ticket := models.Ticket{
		ID:          uuid.NewString(),
		Title:       req.Title,
		Description: req.Description,
		Severity:    req.Severity,
		Status:      models.StatusOpen,
		OwnerID:     userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := repository.CreateTicket(&ticket); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating ticket"})
	}

	logAction(c, userID, "CREATE_TICKET", "ticket", ticket.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": ticket})
}

// GET /api/tickets
func ListTickets(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	userRole := c.Locals("userRole").(string)

	severity := c.Query("severity")
	status := c.Query("status")

	var (
		tickets []models.Ticket
		err     error
	)

	if userRole == "manager" {
		tickets, err = repository.SearchTickets(severity, status)
	} else {
		tickets, err = repository.GetTicketsByOwner(userID)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching tickets"})
	}

	return c.JSON(fiber.Map{"success": true, "data": tickets})
}

// GET /api/tickets/:id
func GetTicket(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	userRole := c.Locals("userRole").(string)

	ticket, err := repository.GetTicketByID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ticket not found"})
	}

	if userRole != "manager" && ticket.OwnerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
	}

	logAction(c, userID, "VIEW_TICKET", "ticket", ticket.ID)

	return c.JSON(fiber.Map{"success": true, "data": ticket})
}

// PUT /api/tickets/:id
func UpdateTicket(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	userRole := c.Locals("userRole").(string)

	ticket, err := repository.GetTicketByID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ticket not found"})
	}

	if userRole != "manager" && ticket.OwnerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
	}

	var req UpdateTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Error parsing"})
	}

	ticket.Title = req.Title
	ticket.Description = req.Description
	ticket.Severity = req.Severity
	ticket.UpdatedAt = time.Now()

	if err := repository.UpdateTicket(ticket); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating ticket"})
	}

	logAction(c, userID, "UPDATE_TICKET", "ticket", ticket.ID)

	return c.JSON(fiber.Map{"success": true, "data": ticket})
}

// DELETE /api/tickets/:id
func DeleteTicket(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	userRole := c.Locals("userRole").(string)

	ticket, err := repository.GetTicketByID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ticket not found"})
	}

	if userRole != "manager" && ticket.OwnerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
	}

	if err := repository.DeleteTicket(ticket.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting ticket"})
	}

	logAction(c, userID, "DELETE_TICKET", "ticket", ticket.ID)

	return c.JSON(fiber.Map{"success": true, "message": "Ticket deleted"})
}

// PATCH /api/tickets/:id/status  (manager only)
func ChangeStatus(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ticket, err := repository.GetTicketByID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ticket not found"})
	}

	var req ChangeStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Error parsing"})
	}

	ticket.Status = req.Status
	ticket.UpdatedAt = time.Now()

	if err := repository.UpdateTicket(ticket); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating status"})
	}

	logAction(c, userID, "CHANGE_STATUS", "ticket", ticket.ID)

	return c.JSON(fiber.Map{"success": true, "data": ticket})
}
