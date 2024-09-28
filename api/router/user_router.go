package api

import (
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeUserRoutes(r *gin.Engine, userHandler *handler.UserHandler) {

	r.Use(middleware.SecurityHeadersMiddleware())
	// Create the 'api' group
	api := r.Group("/api")
	{
		// Create the 'auth' subgroup under the 'api' group
		authGroup := api.Group("/auth")
		{
			authGroup.Use(middleware.JwtAuthMiddleware())

			// Define routes in the 'auth' subgroup
			authGroup.POST("/users", userHandler.CreateUser)
			authGroup.GET("/users/:username", userHandler.GetUserByUsername)
			authGroup.GET("/users", userHandler.ListUsers)
			authGroup.PUT("/users", userHandler.UpdateUser)
			authGroup.DELETE("/users/:id", userHandler.DeleteUser)

		}

		api.POST("/login", userHandler.Login)
	}
}
