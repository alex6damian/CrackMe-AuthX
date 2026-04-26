package main

import (
	"log"

	"github.com/alex6damian/CrackMe-AuthX/api"
	"github.com/alex6damian/CrackMe-AuthX/internal/db"
	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

func main() {
	db.Init()

	engine := html.New("./views", ".html") // views/*html

	app := fiber.New(fiber.Config{
		AppName: "CrackMe by AuthX",
		Views:   engine,
	})

	app.Use(logger.New())
	app.Static("/static", "./static")

	// HTML pages
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/login")
	})
	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{"Title": "Login"})
	})
	app.Get("/register", func(c *fiber.Ctx) error {
		return c.Render("register", fiber.Map{"Title": "Register"})
	})
	app.Get("/profile", api.AuthMiddleware, func(c *fiber.Ctx) error {
		return c.Render("profile", fiber.Map{"Title": "Profile"})
	})
	app.Get("/forgot-password", func(c *fiber.Ctx) error {
		return c.Render("forgot_password", fiber.Map{"Title": "Forgot password"})
	})
	app.Get("/reset-password", func(c *fiber.Ctx) error {
		// Set token from query
		return c.Render("reset_password", fiber.Map{
			"Title": "Reset password",
			"Token": c.Query("token"),
		})
	})

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":   "ok",
			"database": "connected",
			"service":  "crackme-api",
			"version":  "1.0.0",
		})
	})

	// API routes
	v1 := app.Group("/api")
	v1.Post("/register", api.Register)
	v1.Post("/login", api.Login)
	v1.Post("/logout", api.AuthMiddleware, api.Logout)
	v1.Post("/forgot-password", api.ForgotPassword)
	v1.Post("/reset-password", api.ResetPassword)
	v1.Get("/profile", api.AuthMiddleware, api.ViewProfile)

	// Ticket routes
	tickets := v1.Group("/tickets", api.AuthMiddleware)
	tickets.Post("/", api.CreateTicket)
	tickets.Get("/", api.ListTickets)
	tickets.Get("/:id", api.GetTicket)
	tickets.Put("/:id", api.UpdateTicket)
	tickets.Delete("/:id", api.DeleteTicket)
	tickets.Patch("/:id/status", api.RoleMiddleware("manager"), api.ChangeStatus)

	// Audit routes (manager only)
	v1.Get("/audit", api.AuthMiddleware, api.RoleMiddleware("manager"), api.ListAuditLogs)

	log.Println("Server started on: 8080")
	log.Fatal(app.Listen(":8080"))

}
