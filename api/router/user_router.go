package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeUserRoutes(r *gin.Engine, userHandler *handler.UserHandler) {

	// API routes group
	api := r.Group(constants.APIBasePath)
	{
		users := api.Group("/users")
		{
			// Public routes - OAuth2 Authentication
			authGroup := users.Group("/auth")
			{
				authGroup.GET("/google", userHandler.LoginWithGoogle)
				authGroup.GET("/google/callback", userHandler.GoogleCallback)
			}

			// Protected routes - Require JWT authentication
			users.Use(middleware.JwtAuthMiddleware())
			{
				// User read operations
				users.GET("", userHandler.GetPaginatedUsers)
				users.GET("/:username", userHandler.GetUserByUsername)
			}

			// Admin routes - Require RBAC admin role
			adminGroup := users.Group("")
			adminGroup.Use(middleware.RBACMiddleware(constants.RoleAdmin))
			{
				createGroup := adminGroup.Group("")
				{
					createGroup.POST("", userHandler.CreateUser)
				}

				// User modifications - Admin only
				adminGroup.PUT("/:id", userHandler.UpdateUser)
				adminGroup.DELETE("/:id", userHandler.DeleteUser)
			}
		}
	}
}
