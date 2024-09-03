package helper

import (
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
