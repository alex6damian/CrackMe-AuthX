package api

import (
	"time"

	"github.com/alex6damian/CrackMe-AuthX/internal/db"
	"github.com/alex6damian/CrackMe-AuthX/internal/models"
	"github.com/alex6damian/CrackMe-AuthX/internal/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthRequest struct {
	Email    string `json:"email" validate:"email,required,min=3,max=30"`
	Password string `json:"password" validate:"required,strong_password"`
}

type AuthResponse struct {
	Email     string    `json:"user"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"time"`
}

type LogoutRequest struct {
	Email string `json:"user" validate:"required"`
	Token string `json:"token" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"email,required,min=3,max=30"`
}

type ForgotPasswordResponse struct {
	Token string `json:"token"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,strong_password"`
}

// Register handler - POST /api/register
func Register(c *fiber.Ctx) error {
	var req AuthRequest

	// Parse and validate request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing"})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check if exists
	var existingUser models.User
	if err := db.DB.Where("email=?", req.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Wrong credentials"})
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Hashing failed",
		})
	}

	// Create user
	user := models.User{
		Email:         req.Email,
		Password_hash: hashedPassword,
	}

	// Insert user to DB
	if err := db.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating user"})
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating token"})
	}

	// Response
	response := AuthResponse{
		Email:     user.Email,
		Token:     token,
		CreatedAt: user.CreatedAt,
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Lax",
		// Secure:   true,
	})

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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing"})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Find user by email
	var user models.User
	if err := db.DB.Where("email=?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}

	// Check if account is locked
	if user.Locked {
		return c.Status(fiber.StatusLocked).JSON(fiber.Map{"error": "Account locked due to too many failed attempts"})
	}

	// Check password
	if !utils.CheckPassword(user.Password_hash, req.Password) {
		user.LoginAttempts++
		if user.LoginAttempts >= 5 {
			user.Locked = true
		}
		db.DB.Save(&user)
		logAction(c, user.ID, "FAILED_LOGIN", "auth", "")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Reset attempts on successful login
	user.LoginAttempts = 0
	db.DB.Save(&user)

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating token"})
	}

	// Response
	response := AuthResponse{
		Email: user.Email,
		Token: token,
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Lax",
		// Secure:   true,
	})

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

// POST api/forgot-password
func ForgotPassword(c *fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing"})
	}

	// Search DB for user
	var user models.User
	if err := db.DB.Where("email=?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}

	token := utils.GenerateSecureResetToken()
	user.Reset_password = &token

	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error saving token"})
	}

	response := ForgotPasswordResponse{
		Token: *user.Reset_password,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// Password Reset handler - POST /api/reset-password
func ResetPassword(c *fiber.Ctx) error {
	var req ResetPasswordRequest
	// Parse and validate request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing"})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var user models.User
	if err := db.DB.Where("reset_password = ?", req.Token).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}

	user.Password_hash = req.Password
	user.Reset_password = nil

	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to reset password"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password has been reset",
	})
}

// View profile - GET /api/profile
func ViewProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userID = userID.(uint)

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}

	return c.JSON(fiber.Map{
		"message": "Profile",
		"data": fiber.Map{
			"email":      user.Email,
			"role":       user.Role,
			"created_at": user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
	})
}
