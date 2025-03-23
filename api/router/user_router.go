package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeUserRoutes(r *gin.Engine, userHandler *handler.UserHandler) {
	r.Use(middleware.SecurityHeadersMiddleware())

	api := r.Group(constants.APIBasePath)
	{
		// OAuth2 endpoints
		api.GET("/auth/google", userHandler.LoginWithGoogle)
		api.GET("/auth/google/callback", userHandler.GoogleCallback)

		users := api.Group("/users")
		{
			users.Use(middleware.JwtAuthMiddleware())

			users.GET("", userHandler.ListUsers)
			users.GET("/:username", userHandler.GetUserByUsername)

			users.Use(middleware.RBACMiddleware(constants.RoleAdmin))
			{
				users.POST("", userHandler.CreateUser)
				users.PUT("/:username", userHandler.UpdateUser)
				users.DELETE("/:username", userHandler.DeleteUser)
			}
		}
	}
}
