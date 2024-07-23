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

func HandleError(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.Code, appErr)
	} else {
		c.JSON(http.StatusInternalServerError, NewAppError(http.StatusInternalServerError, "Internal Server Error", err.Error()))
	}
}

func HandleSuccess(c *gin.Context, code int, data interface{}, detail string) {
	c.JSON(code, NewAppSuccess(code, data, detail))
}
