package repository

import (
	"github.com/alex6damian/CrackMe-AuthX/internal/db"
	"github.com/alex6damian/CrackMe-AuthX/internal/models"
)

func CreateAuditLog(entry *models.AuditLog) error {
	return db.DB.Create(entry).Error
}

func GetAllAuditLogs() ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := db.DB.Order("timestamp desc").Find(&logs).Error
	return logs, err
}

func GetAuditLogsByUser(userID uint) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := db.DB.Where("user_id = ?", userID).Order("timestamp desc").Find(&logs).Error
	return logs, err
}
