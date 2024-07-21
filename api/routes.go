package api

import (
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(r *gin.Engine, userHandler *handler.UserHandler) {

	r.Use(middleware.SecurityHeadersMiddleware())
	// Create the 'api' group
	api := r.Group("/api")
	{
		// Create the 'auth' subgroup under the 'api' group
		authGroup := api.Group("/auth")
		authGroup.Use(middleware.AuthMiddleware())

		// Define routes in the 'auth' subgroup
		authGroup.GET("/users", userHandler.ListUsers)
		authGroup.POST("/users", userHandler.CreateUser)
	}
}
