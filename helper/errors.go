package helper

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func NewAppError(code int, message, detail string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

func (e *AppError) Error() string {
	return e.Message
}

func HandleError(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.Code, appErr)
	} else {
		c.JSON(http.StatusInternalServerError, NewAppError(http.StatusInternalServerError, "Internal Server Error", err.Error()))
	}
}
