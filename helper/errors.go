package helper

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

type AppSuccess struct {
	Code   int         `json:"code"`
	Data   interface{} `json:"data"`
	Detail string      `json:"detail,omitempty"`
}

func NewAppError(code int, message, detail string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Detail:  detail,
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

func HandleError(c *gin.Context, appErr *AppError, exposeToUser bool) {
	if exposeToUser {
		// Respuesta genérica o mensaje seguro para el usuario
		c.JSON(appErr.Code, appErr)
	} else {
		// Mensaje genérico al usuario, sin detalles tecnicos
		c.JSON(appErr.Code, gin.H{"code": appErr.Code, "error": appErr.Message})
	}
}

func HandleSuccess(c *gin.Context, code int, data interface{}, detail string) {
	c.JSON(code, NewAppSuccess(code, data, detail))
}

func HandleValidationError(c *gin.Context, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errorMessages := make(map[string]string)
		for _, fieldError := range validationErrors {
			errorMessages[fieldError.Field()] = fieldError.Error()
		}
		errorMessagesJSON, _ := json.Marshal(errorMessages)
		HandleError(c, NewAppError(http.StatusBadRequest, "Invalid input", string(errorMessagesJSON)), true)
	} else {
		HandleError(c, NewAppError(http.StatusBadRequest, "Invalid input", err.Error()), true)
	}
}

func HandleGormError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		HandleError(c, NewAppError(http.StatusNotFound, "Resource not found", err.Error()), true)
	case errors.Is(err, gorm.ErrDuplicatedKey):
		HandleError(c, NewAppError(http.StatusConflict, "Duplicated key", err.Error()), true)
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		HandleError(c, NewAppError(http.StatusBadRequest, "Foreign key constraint violated", err.Error()), true)
	case errors.Is(err, gorm.ErrInvalidTransaction):
		HandleError(c, NewAppError(http.StatusBadRequest, "Invalid transaction", err.Error()), true)
	case errors.Is(err, gorm.ErrNotImplemented):
		HandleError(c, NewAppError(http.StatusNotImplemented, "Not implemented", err.Error()), true)
	case errors.Is(err, gorm.ErrMissingWhereClause):
		HandleError(c, NewAppError(http.StatusBadRequest, "Missing WHERE clause", err.Error()), true)
	case errors.Is(err, gorm.ErrUnsupportedRelation):
		HandleError(c, NewAppError(http.StatusBadRequest, "Unsupported relation", err.Error()), true)
	case errors.Is(err, gorm.ErrPrimaryKeyRequired):
		HandleError(c, NewAppError(http.StatusBadRequest, "Primary key required", err.Error()), true)
	case errors.Is(err, gorm.ErrModelValueRequired):
		HandleError(c, NewAppError(http.StatusBadRequest, "Model value required", err.Error()), true)
	case errors.Is(err, gorm.ErrModelAccessibleFieldsRequired):
		HandleError(c, NewAppError(http.StatusBadRequest, "Model accessible fields required", err.Error()), true)
	case errors.Is(err, gorm.ErrSubQueryRequired):
		HandleError(c, NewAppError(http.StatusBadRequest, "Sub query required", err.Error()), true)
	case errors.Is(err, gorm.ErrInvalidData):
		HandleError(c, NewAppError(http.StatusBadRequest, "Invalid data", err.Error()), true)
	case errors.Is(err, gorm.ErrUnsupportedDriver):
		HandleError(c, NewAppError(http.StatusBadRequest, "Unsupported driver", err.Error()), true)
	case errors.Is(err, gorm.ErrRegistered):
		HandleError(c, NewAppError(http.StatusConflict, "Already registered", err.Error()), true)
	case errors.Is(err, gorm.ErrInvalidField):
		HandleError(c, NewAppError(http.StatusBadRequest, "Invalid field", err.Error()), true)
	case errors.Is(err, gorm.ErrEmptySlice):
		HandleError(c, NewAppError(http.StatusBadRequest, "Empty slice found", err.Error()), true)
	case errors.Is(err, gorm.ErrDryRunModeUnsupported):
		HandleError(c, NewAppError(http.StatusBadRequest, "Dry run mode unsupported", err.Error()), true)
	case errors.Is(err, gorm.ErrInvalidDB):
		HandleError(c, NewAppError(http.StatusBadRequest, "Invalid DB", err.Error()), true)
	case errors.Is(err, gorm.ErrInvalidValue):
		HandleError(c, NewAppError(http.StatusBadRequest, "Invalid value", err.Error()), true)
	case errors.Is(err, gorm.ErrInvalidValueOfLength):
		HandleError(c, NewAppError(http.StatusBadRequest, "Invalid value length", err.Error()), true)
	case errors.Is(err, gorm.ErrPreloadNotAllowed):
		HandleError(c, NewAppError(http.StatusBadRequest, "Preload not allowed", err.Error()), true)
	case errors.Is(err, gorm.ErrCheckConstraintViolated):
		HandleError(c, NewAppError(http.StatusBadRequest, "Check constraint violated", err.Error()), true)
	default:
		HandleError(c, NewAppError(http.StatusInternalServerError, "Internal server error", err.Error()), true)
	}
}
