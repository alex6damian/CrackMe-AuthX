package models

import "time"

type User struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Email          string    `gorm:"uniqueIndex;not null" json:"email"`
	Password_hash  string    `gorm:"not null" json:"-"`
	Reset_password *string   `gorm:"uniqueIndex" json:"-"`
	Role           string    `gorm:"default: user" json:"role"`
	CreatedAt      time.Time `json:"created_at"`
	Locked         bool      `gorm:"default: false" json:"locked"`
}
