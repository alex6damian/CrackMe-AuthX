package api

import (
	"fmt"

	"github.com/alex6damian/CrackMe-AuthX/internal/db"
	"github.com/alex6damian/CrackMe-AuthX/internal/models"
	"github.com/alex6damian/CrackMe-AuthX/internal/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthRequest struct {
	Email    string `json:"email" validate:"required, min:3, max 30"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Email string `json:"user"`
	Token string `json:"token"`
}

// Register handler - POST /api/register
func Register(c *fiber.Ctx) error {
	var req AuthRequest

	// Parse and validate request
	if err := c.BodyParser(&req); err != nil {
		return fmt.Errorf("Error parsing: %v", err)
	}

	// Check if exists
	var existingUser models.User
	if err := db.DB.Where("email=?", req.Email).First(&existingUser).Error; err == nil {
		return fmt.Errorf("user already existing: %v", err)
	}

	// TODO: hash passs

	// Create user
	user := models.User{
		Email:         req.Email,
		Password_hash: req.Password,
	}

	// Insert user to DB
	if err := db.DB.Create(&user).Error; err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return fmt.Errorf("Error creating JWT token: %v", err)
	}

	// Response
	response := AuthResponse{
		Email: user.Email,
		Token: token,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// Login handler - POST /api/login
func Login(c *fiber.Ctx) error {
	var req AuthRequest

	// Parse and validate request
	if err := c.BodyParser(&req); err != nil {
		return fmt.Errorf("Error parsing: %v", err)
	}

	// raw sql for injection

	// Find user by email
	var user models.User
	if err := db.DB.Where("email=?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("User not found: %v", err)
		}
		return fmt.Errorf("DB error: %v", err)
	}

	// Check password
	if user.Password_hash != req.Password {
		return fmt.Errorf("Wrong password!")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return fmt.Errorf("Error generating JWT: %v", err)
	}

	// Response
	response := AuthResponse{
		Email: user.Email,
		Token: token,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// Lougout handler - POST /api/logout
func Logout(c *fiber.Ctx) error {
	/*
		Pentru securizare, invalidare token instant
	*/

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "You have been successfully logged out!",
	})
}

// Password Reset handler - POST /api/reset
func ResetPassword(c *fiber.Ctx) error {
	return fmt.Errorf("In dezvoltare")
}
