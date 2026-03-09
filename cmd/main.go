package main

import (
	"log"

	"github.com/alex6damian/CrackMe-AuthX/internal/db"
	"github.com/gofiber/fiber/v2"
)

func main() {
	db.Init()

	app := fiber.New(fiber.Config{
		AppName: "CrackMe by AuthX",
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":   "ok",
			"database": "connected",
			"service":  "crackme-api",
			"version":  "1.0.0",
		})
	})

	log.Println("Server started on: 8080")
	log.Fatal(app.Listen(":8080"))

}
