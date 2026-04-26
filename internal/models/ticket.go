package models

import "time"

type Severity string
type TicketStatus string

const (
	SeverityLow  Severity = "LOW"
	SeverityMed  Severity = "MED"
	SeverityHigh Severity = "HIGH"
)

const (
	StatusOpen       TicketStatus = "OPEN"
	StatusInProgress TicketStatus = "IN_PROGRESS"
	StatusResolved   TicketStatus = "RESOLVED"
)

type Ticket struct {
	ID          string       `gorm:"primaryKey;type:text" json:"id"`
	Title       string       `gorm:"not null" json:"title"`
	Description string       `gorm:"type:text" json:"description"`
	Severity    Severity     `gorm:"type:text;not null" json:"severity"`
	Status      TicketStatus `gorm:"type:text;default:'OPEN'" json:"status"`
	OwnerID     uint         `gorm:"not null" json:"owner_id"`
	Owner       User         `gorm:"foreignKey:OwnerID" json:"-"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
