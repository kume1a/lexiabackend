package shared

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	UsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

	lowercaseRegex = regexp.MustCompile(`[a-z]`)
	uppercaseRegex = regexp.MustCompile(`[A-Z]`)
	digitRegex     = regexp.MustCompile(`\d`)
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   any    `json:"value,omitempty"`
}

type ValidationErrorResponse struct {
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors"`
}

func ValidateUsername(fl validator.FieldLevel) bool {
	return UsernameRegex.MatchString(fl.Field().String())
}

func ValidateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}
	return lowercaseRegex.MatchString(password) &&
		uppercaseRegex.MatchString(password) &&
		digitRegex.MatchString(password)
}

func RegisterCustomValidations(v *validator.Validate) {
	v.RegisterValidation("username", ValidateUsername)
	v.RegisterValidation("strong_password", ValidateStrongPassword)
}

func ValidateStruct(s any) *ValidationErrorResponse {
	validate := validator.New()
	RegisterCustomValidations(validate)

	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors []ValidationError

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			validationErrors = append(validationErrors, ValidationError{
				Field:   strings.ToLower(e.Field()),
				Message: getValidationMessage(e),
				Value:   e.Value(),
			})
		}
	}

	return &ValidationErrorResponse{
		Message: "Validation failed",
		Errors:  validationErrors,
	}
}

func BindAndValidate(c *gin.Context, obj any) *ValidationErrorResponse {
	if err := c.ShouldBindJSON(obj); err != nil {
		return handleBindingError(err)
	}
	return nil
}

func handleBindingError(err error) *ValidationErrorResponse {
	var validationErrors []ValidationError

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			validationErrors = append(validationErrors, ValidationError{
				Field:   strings.ToLower(e.Field()),
				Message: getValidationMessage(e),
				Value:   e.Value(),
			})
		}
	} else {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "body",
			Message: "Invalid JSON format or request body",
		})
	}

	return &ValidationErrorResponse{
		Message: "Validation failed",
		Errors:  validationErrors,
	}
}

func getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return "Must be a valid email address"
	case "min":
		if e.Kind().String() == "string" {
			return fmt.Sprintf("Must be at least %s characters long", e.Param())
		}
		return fmt.Sprintf("Must be at least %s", e.Param())
	case "max":
		if e.Kind().String() == "string" {
			return fmt.Sprintf("Must be no more than %s characters long", e.Param())
		}
		return fmt.Sprintf("Must be no more than %s", e.Param())
	case "alphanum":
		return "Must contain only letters and numbers"
	case "alpha":
		return "Must contain only letters"
	case "numeric":
		return "Must be a valid number"
	case "url":
		return "Must be a valid URL"
	case "uuid":
		return "Must be a valid UUID"
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", e.Param())
	case "len":
		return fmt.Sprintf("Must be exactly %s characters long", e.Param())
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", e.Param())
	case "lte":
		return fmt.Sprintf("Must be less than or equal to %s", e.Param())
	case "gt":
		return fmt.Sprintf("Must be greater than %s", e.Param())
	case "lt":
		return fmt.Sprintf("Must be less than %s", e.Param())
	case "username":
		return "Username can only contain letters, numbers, underscores, and hyphens"
	case "strong_password":
		return "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, and one number"
	default:
		return fmt.Sprintf("Invalid %s", e.Field())
	}
}

func ResValidationError(c *gin.Context, validationError *ValidationErrorResponse) {
	c.JSON(400, validationError)
}
