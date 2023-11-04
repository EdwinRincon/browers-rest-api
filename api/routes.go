package api

import (
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(r *gin.Engine, userHandler *handler.UserHandler) {

	authGroup := r.Group("/auth")
	authGroup.Use(middleware.AuthMiddleware())

	authGroup.GET("/users", userHandler.ListUsers)
	authGroup.POST("/users", userHandler.CreateUser)
}
