package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeUserRoutes(r *gin.Engine, userHandler *handler.UserHandler) {
	// Global middleware for all routes
	r.Use(middleware.SecurityHeadersMiddleware())

	// API routes group
	api := r.Group(constants.APIBasePath)
	{
		users := api.Group("/users")
		{
			// Public routes - OAuth2 Authentication
			// Rate limited to prevent auth abuse
			authGroup := users.Group("/auth")
			authGroup.Use(middleware.RateLimitAuth())
			{
				authGroup.GET("/google", userHandler.LoginWithGoogle)
				authGroup.GET("/google/callback", userHandler.GoogleCallback)
			}

			// Protected routes - Require JWT authentication
			users.Use(middleware.JwtAuthMiddleware())
			{
				// User read operations
				users.GET("", userHandler.ListUsers)                   // List all users
				users.GET("/:username", userHandler.GetUserByUsername) // Get specific user
			}

			// Admin routes - Require RBAC admin role
			adminGroup := users.Group("")
			adminGroup.Use(middleware.RBACMiddleware(constants.RoleAdmin))
			{
				// User creation - Special rate limiting
				createGroup := adminGroup.Group("")
				createGroup.Use(middleware.RateLimitNewAccounts())
				{
					createGroup.POST("", userHandler.CreateUser) // Create new user
				}

				// User modifications - Admin only
				adminGroup.PUT("/:id", userHandler.UpdateUser)    // Update user
				adminGroup.DELETE("/:id", userHandler.DeleteUser) // Delete user
			}
		}
	}
}
