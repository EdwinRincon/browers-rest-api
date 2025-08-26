package helper

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/EdwinRincon/browersfc-api/pkg/logger"
	"github.com/EdwinRincon/browersfc-api/pkg/validation"
	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code        int               `json:"code"`
	Message     string            `json:"message"`
	Detail      string            `json:"detail,omitempty"`
	Validation  map[string]string `json:"validation,omitempty"`
	Field       string            `json:"field,omitempty"`
	SafeForUser bool              `json:"-"`
	cause       error             `json:"-"` // Internal error for logging, never exposed to client
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
		withCause(err),
	)
}

// ValidationError creates a user-safe validation error with structured field errors
func ValidationError(validation map[string]string) *AppError {
	return newAppError(
		http.StatusBadRequest,
		"Validation failed",
		withDetail("One or more fields are invalid"),
		withValidation(validation),
		safe(),
	)
}

// ProcessValidationError is a convenience function to handle validator errors
// It extracts validation errors if present, or returns a generic bad request error otherwise
func ProcessValidationError(err error, field, defaultMessage string) *AppError {
	validationErrs := validation.ExtractValidationErrors(err)
	if len(validationErrs) > 0 {
		return ValidationError(validationErrs)
	}
	return newAppError(
		http.StatusBadRequest,
		defaultMessage,
		withField(field),
		withDetail(fmt.Sprintf("The %s field has an invalid value", field)),
		withCause(err), // <--- attach the raw Go error
		safe(),
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

// withCause adds the original error as the cause (for internal use only)
func withCause(err error) func(*AppError) {
	return func(e *AppError) {
		e.cause = err
	}
}

// withValidation adds validation errors to the error
func withValidation(validation map[string]string) func(*AppError) {
	return func(e *AppError) {
		e.Validation = validation
	}
}

func (e *AppError) Error() string {
	if e.cause != nil {
		return e.cause.Error()
	}
	return e.Message
}

// Unwrap implements the errors.Unwrap interface
// It allows errors.Is and errors.As to work with wrapped errors
func (e *AppError) Unwrap() error {
	return e.cause
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

	if e.cause != nil {
		attrs = append(attrs, "error_cause", e.cause.Error())
	}

	// Add error attributes to logger
	return l.With(attrs...)
}

// ToResponse returns a map that can be safely serialized to JSON for client responses
// It ensures consistent response shape and prevents leaking sensitive details
func (e *AppError) ToResponse() map[string]interface{} {
	response := map[string]interface{}{
		"code": e.Code,
	}

	// If the error is not safe for users, use generic messages
	if !e.SafeForUser {
		response["message"] = "An error occurred"
		response["detail"] = "Please try again or contact support if the problem persists"
		return response
	}

	// Otherwise include all safe fields
	response["message"] = e.Message

	if e.Detail != "" {
		response["detail"] = e.Detail
	}

	if e.Field != "" {
		response["field"] = e.Field
	}

	if len(e.Validation) > 0 {
		response["validation"] = e.Validation
	}

	return response
}

// RespondWithError sends an error response to the client
func RespondWithError(c *gin.Context, err *AppError) {
	// Store the error in the context for later logging by the middleware
	logger.StoreErrorForLogging(c, err)

	// Use the ToResponse method to ensure a consistent response shape
	// This will automatically handle SafeForUser flag
	c.JSON(err.Code, err.ToResponse())

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
