package helper

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
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

func HandleSuccess(c *gin.Context, code int, data interface{}, detail string) {
	c.JSON(code, NewAppSuccess(code, data, detail))
}
