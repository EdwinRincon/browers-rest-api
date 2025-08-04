package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type AppError struct {
	Code        int               `json:"code"`
	Message     string            `json:"message"`
	Detail      string            `json:"detail,omitempty"`
	Validation  map[string]string `json:"validation,omitempty"`
	Field       string            `json:"field,omitempty"`
	SafeForUser bool              `json:"-"`
}

type AppSuccess struct {
	Code   int         `json:"code"`
	Data   interface{} `json:"data"`
	Detail string      `json:"detail,omitempty"`
}

// newAppError creates a new application error
func newAppError(code int, message string, opts ...func(*AppError)) *AppError {
	err := &AppError{
		Code:    code,
		Message: message,
	}

	for _, opt := range opts {
		opt(err)
	}

	return err
}

// BadRequest creates a user-safe 400 error
func BadRequest(field, message string) *AppError {
	return newAppError(
		http.StatusBadRequest,
		message,
		WithField(field),
		WithDetail(fmt.Sprintf("The %s field has an invalid value", field)),
		Safe(),
	)
}

// Unauthorized creates a user-safe 401 error
func Unauthorized(detail string) *AppError {
	return newAppError(
		http.StatusUnauthorized,
		"Unauthorized access",
		WithDetail(detail),
		Safe(),
	)
}

// StatusForbidden creates a user-safe 403 error
func StatusForbidden(detail string) *AppError {
	return newAppError(
		http.StatusForbidden,
		"Forbidden access",
		WithDetail(detail),
		Safe(),
	)
}

// NotFound creates a user-safe 404 error
func NotFound(resource string) *AppError {
	return newAppError(
		http.StatusNotFound,
		fmt.Sprintf("%s not found", resource),
		WithDetail(fmt.Sprintf("The requested %s could not be found", resource)),
		Safe(),
	)
}

// Conflict creates a user-safe 409 error
func Conflict(resource string, detail string) *AppError {
	return newAppError(
		http.StatusConflict,
		fmt.Sprintf("%s already exists", resource),
		WithDetail(detail),
		Safe(),
	)
}

// InternalError creates a 500 error that logs but doesn't expose details
func InternalError(err error) *AppError {
	// Log the actual error
	slog.Error("Internal server error", "error", err)
	return newAppError(
		http.StatusInternalServerError,
		"Internal server error",
		WithDetail("An unexpected error occurred"),
	)
}

// Safe marks an error as safe to expose to users
func Safe() func(*AppError) {
	return func(e *AppError) {
		e.SafeForUser = true
	}
}

// WithDetail adds detail to the error
func WithDetail(detail string) func(*AppError) {
	return func(e *AppError) {
		e.Detail = detail
	}
}

// WithValidation adds validation errors
func WithValidation(validationErrors map[string]string) func(*AppError) {
	return func(e *AppError) {
		e.Validation = validationErrors
	}
}

// WithField adds field information
func WithField(field string) func(*AppError) {
	return func(e *AppError) {
		e.Field = field
	}
}

func NewAppSuccess(code int, data interface{}, detail string) *AppSuccess {
	return &AppSuccess{
		Code:   code,
		Data:   data,
		Detail: detail,
	}
}

func (e *AppError) Error() string {
	return e.Message
}

// RespondWithError sends an error response to the client
func RespondWithError(c *gin.Context, err *AppError) {
	if err.SafeForUser {
		c.JSON(err.Code, err)
	} else {
		// For unsafe errors, only expose the status code and a generic message
		c.JSON(err.Code, gin.H{
			"code":    err.Code,
			"message": "An error occurred",
			"detail":  "Please try again or contact support if the problem persists",
		})
	}
}

// handleError is maintained for backward compatibility
func handleError(c *gin.Context, appErr *AppError, exposeToUser bool) {
	if exposeToUser {
		appErr.SafeForUser = true
	}
	RespondWithError(c, appErr)
}

func HandleSuccess(c *gin.Context, code int, data interface{}, detail string) {
	c.JSON(code, NewAppSuccess(code, data, detail))
}

func HandleValidationError(c *gin.Context, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errorMessages := make(map[string]string)
		var firstField string
		for _, fieldError := range validationErrors {
			if firstField == "" {
				firstField = fieldError.Field()
			}
			switch fieldError.Tag() {
			case "required":
				errorMessages[fieldError.Field()] = "This field is required"
			case "min":
				errorMessages[fieldError.Field()] = fmt.Sprintf("Minimum value is %s", fieldError.Param())
			case "max":
				errorMessages[fieldError.Field()] = fmt.Sprintf("Maximum value is %s", fieldError.Param())
			case "gte":
				errorMessages[fieldError.Field()] = fmt.Sprintf("Must be greater than or equal to %s", fieldError.Param())
			case "lte":
				errorMessages[fieldError.Field()] = fmt.Sprintf("Must be less than or equal to %s", fieldError.Param())
			case "email":
				errorMessages[fieldError.Field()] = "Must be a valid email address"
			case "url":
				errorMessages[fieldError.Field()] = "Must be a valid URL"
			case "safe_email":
				errorMessages[fieldError.Field()] = "Must be a valid and secure email address"
			case "allowed_domain":
				errorMessages[fieldError.Field()] = "Email domain is not allowed. Please use a supported email provider"
			case "alphanum":
				errorMessages[fieldError.Field()] = "Must contain only letters and numbers"
			case "oneof":
				errorMessages[fieldError.Field()] = fmt.Sprintf("Must be one of: %s", fieldError.Param())
			case "len":
				errorMessages[fieldError.Field()] = fmt.Sprintf("Must be exactly %s characters long", fieldError.Param())
			case "numeric":
				errorMessages[fieldError.Field()] = "Must contain only numbers"
			case "uuid":
				errorMessages[fieldError.Field()] = "Must be a valid UUID"
			case "datetime":
				errorMessages[fieldError.Field()] = "Must be a valid date and time in format YYYY-MM-DDTHH:MM:SSZ"
			case "date":
				errorMessages[fieldError.Field()] = "Must be a valid date in format YYYY-MM-DD"
			case "time":
				errorMessages[fieldError.Field()] = "Must be a valid time in format HH:MM"
			default:
				// Log unexpected validation tags for monitoring
				slog.Debug("Unhandled validation tag",
					"tag", fieldError.Tag(),
					"field", fieldError.Field(),
					"param", fieldError.Param())
				errorMessages[fieldError.Field()] = fieldError.Error()
			}
		}

		appError := newAppError(
			http.StatusBadRequest,
			"Validation failed",
			WithValidation(errorMessages),
			WithField(firstField),
		)
		handleError(c, appError, true)
	} else if jsonErr, ok := err.(*json.UnmarshalTypeError); ok {
		message := fmt.Sprintf("Invalid value type for field '%s'. Expected %s, got %s",
			jsonErr.Field, jsonErr.Type, jsonErr.Value)
		appError := newAppError(
			http.StatusBadRequest,
			"Invalid input type",
			WithDetail(message),
			WithField(jsonErr.Field),
		)
		handleError(c, appError, true)
	} else {
		appError := newAppError(
			http.StatusBadRequest,
			"Invalid input format",
			WithDetail("Please check the input format and try again"),
		)
		handleError(c, appError, true)
	}
}

// isValidationError checks if the error is a GORM validation error
func isValidationError(err error) bool {
	return errors.Is(err, gorm.ErrInvalidData) ||
		errors.Is(err, gorm.ErrInvalidField) ||
		errors.Is(err, gorm.ErrInvalidValue) ||
		errors.Is(err, gorm.ErrInvalidValueOfLength)
}

// HandleGormError handles database errors in a consistent way
func HandleGormError(c *gin.Context, err error) {
	var appErr *AppError

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		appErr = NotFound("resource")
	case errors.Is(err, gorm.ErrDuplicatedKey):
		appErr = Conflict("resource", "A record with these details already exists")
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		appErr = BadRequest("", "The referenced resource does not exist")
	case isValidationError(err):
		appErr = BadRequest("", "The provided data is invalid")
	case errors.Is(err, gorm.ErrNotImplemented):
		appErr = newAppError(
			http.StatusNotImplemented,
			"Feature not available",
			WithDetail("This feature is not implemented yet"),
			Safe(),
		)
	default:
		// Log unexpected errors but don't expose details to user
		appErr = InternalError(err)
	}

	RespondWithError(c, appErr)
}
