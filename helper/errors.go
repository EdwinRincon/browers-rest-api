package helper

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/EdwinRincon/browersfc-api/pkg/logger"
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
		withField(field),
		withDetail(fmt.Sprintf("The %s field has an invalid value", field)),
		safe(),
	)
}

// Unauthorized creates a user-safe 401 error
func Unauthorized(detail string) *AppError {
	return newAppError(
		http.StatusUnauthorized,
		"Unauthorized access",
		withDetail(detail),
		safe(),
	)
}

// StatusForbidden creates a user-safe 403 error
func StatusForbidden(detail string) *AppError {
	return newAppError(
		http.StatusForbidden,
		"Forbidden access",
		withDetail(detail),
		safe(),
	)
}

// NotFound creates a user-safe 404 error
func NotFound(resource string) *AppError {
	return newAppError(
		http.StatusNotFound,
		fmt.Sprintf("%s not found", resource),
		withDetail(fmt.Sprintf("The requested %s could not be found", resource)),
		safe(),
	)
}

// Conflict creates a user-safe 409 error
func Conflict(resource string, detail string) *AppError {
	return newAppError(
		http.StatusConflict,
		fmt.Sprintf("%s already exists", resource),
		withDetail(detail),
		safe(),
	)
}

// InternalError creates a 500 error that logs but doesn't expose details
func InternalError(err error) *AppError {
	return newAppError(
		http.StatusInternalServerError,
		"Internal server error",
		withDetail("An unexpected error occurred"),
	)
}

// safe marks an error as safe to expose to users
func safe() func(*AppError) {
	return func(e *AppError) {
		e.SafeForUser = true
	}
}

// withDetail adds detail to the error
func withDetail(detail string) func(*AppError) {
	return func(e *AppError) {
		e.Detail = detail
	}
}

// withField adds field information
func withField(field string) func(*AppError) {
	return func(e *AppError) {
		e.Field = field
	}
}

func (e *AppError) Error() string {
	return e.Message
}

// AddToLog implements the logger.LoggableError interface
// It adds structured error data to the given logger
func (e *AppError) AddToLog(l *slog.Logger) *slog.Logger {
	// Start with required fields
	attrs := []any{
		"error_code", e.Code,
		"error_message", e.Message,
	}

	// Add optional fields if present
	if e.Detail != "" {
		attrs = append(attrs, "error_detail", e.Detail)
	}

	if e.Field != "" {
		attrs = append(attrs, "error_field", e.Field)
	}

	if len(e.Validation) > 0 {
		// Add validation errors as a nested object
		for field, msg := range e.Validation {
			attrs = append(attrs, "validation."+field, msg)
		}
	}

	// Add error attributes to logger
	return l.With(attrs...)
}

// RespondWithError sends an error response to the client
func RespondWithError(c *gin.Context, err *AppError) {
	// Store the error in the context for later logging by the middleware
	logger.StoreErrorForLogging(c, err)

	if err.SafeForUser {
		// safe errors can be shown to users directly
		c.JSON(err.Code, err)
	} else {
		// For unsafe errors, only expose the status code and a generic message
		c.JSON(err.Code, gin.H{
			"code":    err.Code,
			"message": "An error occurred",
			"detail":  "Please try again or contact support if the problem persists",
		})
	}

	// Add error to Gin's error collection for consistency with built-in error handling
	_ = c.Error(err)
}

func newAppSuccess(code int, data interface{}, detail string) *AppSuccess {
	return &AppSuccess{
		Code:   code,
		Data:   data,
		Detail: detail,
	}
}

func HandleSuccess(c *gin.Context, code int, data interface{}, detail string) {
	c.JSON(code, newAppSuccess(code, data, detail))
}
