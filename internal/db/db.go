package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/alex6damian/CrackMe-AuthX/internal/models"
)

var DB *gorm.DB

func Init() {
	DATABASE_URL := os.Getenv("DATABASE_URL")

	if DATABASE_URL == "" {
		DATABASE_URL = "authx.db"
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(DATABASE_URL), &gorm.Config{})
	if err != nil {
		fmt.Printf("Database connection error")
	}

	log.Printf("Database connected")

	RunMigrations()
}

func RunMigrations() {
	log.Println("Running migrations")
	err := DB.AutoMigrate(
		&models.User{},
	)

	if err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}

	log.Println("Migrations completed")
}
