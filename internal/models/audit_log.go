package models

import "time"

type AuditLog struct {
	ID         string    `gorm:"primaryKey;type:text" json:"id"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"-"`
	Action     string    `gorm:"type:text;not null" json:"action"`
	Resource   string    `gorm:"type:text;not null" json:"resource"`
	ResourceID string    `gorm:"type:text" json:"resource_id"`
	TicketID   *string   `gorm:"type:text" json:"ticket_id,omitempty"`
	Ticket     *Ticket   `gorm:"foreignKey:TicketID" json:"-"`
	Timestamp  time.Time `json:"timestamp"`
	IPAddress  string    `gorm:"type:text" json:"ip_address"`
}
