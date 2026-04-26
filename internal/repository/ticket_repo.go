package repository

import (
	"github.com/alex6damian/CrackMe-AuthX/internal/db"
	"github.com/alex6damian/CrackMe-AuthX/internal/models"
)

func CreateTicket(ticket *models.Ticket) error {
	return db.DB.Create(ticket).Error
}

func GetTicketByID(id string) (*models.Ticket, error) {
	var ticket models.Ticket
	err := db.DB.First(&ticket, "id = ?", id).Error
	return &ticket, err
}

func GetAllTickets() ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := db.DB.Find(&tickets).Error
	return tickets, err
}

func GetTicketsByOwner(ownerID uint) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := db.DB.Where("owner_id = ?", ownerID).Find(&tickets).Error
	return tickets, err
}

func UpdateTicket(ticket *models.Ticket) error {
	return db.DB.Save(ticket).Error
}

func DeleteTicket(id string) error {
	return db.DB.Delete(&models.Ticket{}, "id = ?", id).Error
}

func SearchTickets(severity string, status string) ([]models.Ticket, error) {
	var tickets []models.Ticket
	query := db.DB.Model(&models.Ticket{})
	if severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&tickets).Error
	return tickets, err
}
