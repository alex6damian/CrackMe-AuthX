package utils

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// Global validator instance
var validate *validator.Validate

// init runs once when package is imported
func init() {
	validate = validator.New()

	// Register custom validators
	validate.RegisterValidation("strong_password", ValidateStrongPassword)
}

// formatValidationError converts validator.FieldError to user-friendly message
func formatValidationError(e validator.FieldError) error {
	field := e.Field()

	switch e.Tag() {
	case "required":
		return fmt.Errorf("%s is required", field)
	case "email":
		return fmt.Errorf("%s must be a valid email address", field)
	case "min":
		return fmt.Errorf("%s must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Errorf("%s must be at most %s characters", field, e.Param())
	case "strong_password":
		return fmt.Errorf("%s must be at least 8 characters with at least one uppercase letter and one number", field)
	case "oneof":
		return fmt.Errorf("%s must be one of: %s", field, e.Param())
	case "url":
		return fmt.Errorf("%s must be a valid URL", field)
	default:
		return fmt.Errorf("%s is invalid", field)
	}
}

// ValidateStrongPassword validates password strength ("strong_password" tag)
func ValidateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	return hasUpper && hasNumber
}

func ValidateStruct(v any) error {
	if err := validate.Struct(v); err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			return formatValidationError(verrs[0])
		}
		return err
	}
	return nil
}
