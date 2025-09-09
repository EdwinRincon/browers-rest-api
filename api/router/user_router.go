package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/EdwinRincon/browersfc-api/internal/domain/service"
	"github.com/gin-gonic/gin"
)

func InitializeUserRoutes(r *gin.Engine, userHandler *handler.UserHandler, authService *service.AuthenticationDomainService) {
	api := r.Group(constants.APIBasePath)
	{
		// Public routes - OAuth2 Authentication
		authGroup := api.Group("/users/auth")
		{
			authGroup.GET("/google", userHandler.LoginWithGoogle)
			authGroup.GET("/google/callback", userHandler.GoogleCallback)
		}

		// Protected user routes
		users := api.Group("/users")
		users.Use(middleware.JwtAuthMiddleware(authService))
		{
			users.GET("", userHandler.GetPaginatedUsers)
			users.GET("/:username", userHandler.GetUserByUsername)
		}

		// Admin routes
		adminUsers := api.Group("/admin/users")
		adminUsers.Use(middleware.JwtAuthMiddleware(authService), middleware.RBACMiddleware(constants.RoleAdmin))
		{
			adminUsers.POST("", userHandler.CreateUser)
			adminUsers.PUT("/:id", userHandler.UpdateUser)
			adminUsers.DELETE("/:id", userHandler.DeleteUser)
		}
	}
}
