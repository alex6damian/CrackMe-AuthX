package api

import (
	"time"

	"github.com/alex6damian/CrackMe-AuthX/internal/db"
	"github.com/alex6damian/CrackMe-AuthX/internal/models"
	"github.com/alex6damian/CrackMe-AuthX/internal/utils"

	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthRequest struct {
	Email    string `json:"email" validate:"required, min:3, max 30"`
	Password string `json:"password" validate:"required"`
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
	Email string `json:"email" validate:"required, min:3, max 30"`
}

type ForgotPasswordResponse struct {
	Token string `json:"token"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Register handler - POST /api/register
func Register(c *fiber.Ctx) error {
	var req AuthRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing"})
	}

	var existingUser models.User
	if err := db.DB.Where("email=?", req.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "User already existing"})
	}

	user := models.User{
		Email:         req.Email,
		Password_hash: req.Password,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating user"})
	}

	logAction(c, user.ID, "REGISTER", "auth", "")

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating token"})
	}

	response := AuthResponse{
		Email:     user.Email,
		Token:     token,
		CreatedAt: user.CreatedAt,
	}

	c.Cookie(&fiber.Cookie{
		Name:  "token",
		Value: token,
		Path:  "/",
		// HTTPOnly: true,
		SameSite: "Lax",
		// Secure: true,
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// Login handler - POST /api/login
func Login(c *fiber.Ctx) error {
	var req AuthRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing"})
	}

	// raw sql - vulnerable to SQL injection via string concatenation
	var user models.User
	query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", req.Email)
	if err := db.DB.Raw(query).Scan(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	if user.Password_hash != req.Password {
		logAction(c, user.ID, "FAILED_LOGIN", "auth", "")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Wrong password"})
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating token"})
	}

	response := AuthResponse{
		Email: user.Email,
		Token: token,
	}

	c.Cookie(&fiber.Cookie{
		Name:  "token",
		Value: token,
		Path:  "/",
		// HTTPOnly: true,
		SameSite: "Lax",
		// Secure: true,
	})

	logAction(c, user.ID, "LOGIN", "auth", "")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// Logout handler - POST /api/logout
func Logout(c *fiber.Ctx) error {
	userID, _ := c.Locals("userID").(uint)
	logAction(c, userID, "LOGOUT", "auth", "")

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

	var user models.User
	if err := db.DB.Where("email=?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}

	var token string
	if user.Reset_password == nil {
		token = utils.GeneratePredictableResetToken(req.Email)
		user.Reset_password = &token
	}

	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error saving token"})
	}

	logAction(c, user.ID, "FORGOT_PASSWORD", "auth", "")

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
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing"})
	}

	var user models.User
	if err := db.DB.Where("reset_password = ?", req.Token).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}

	user.Password_hash = req.Password

	// De invalidat pe viitor

	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to reset password"})
	}

	logAction(c, user.ID, "RESET_PASSWORD", "auth", "")

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
