package api

import (
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeRoleRoutes(r *gin.Engine, roleHandler *handler.RoleHandler) {

	r.Use(middleware.SecurityHeadersMiddleware())
	// Create the 'api' group
	api := r.Group("/api")
	{
		// Create the 'auth' subgroup under the 'api' group
		authGroup := api.Group("/auth")
		{
			authGroup.Use(middleware.JwtAuthMiddleware())

			// Define routes in the 'auth' subgroup
			authGroup.GET("/role/:id", roleHandler.GetRoleByID)
			authGroup.POST("/role", roleHandler.CreateRole)
			authGroup.PUT("/role/:id", roleHandler.UpdateRole)
			authGroup.DELETE("/role/:id", roleHandler.DeleteRole)
			authGroup.GET("/roles", roleHandler.GetAllRoles)
		}
	}
}
