package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeRoleRoutes(r *gin.Engine, roleHandler *handler.RoleHandler) {
	r.Use(middleware.SecurityHeadersMiddleware())

	api := r.Group(constants.APIBasePath)
	{
		roles := api.Group("/roles")
		{
			roles.Use(middleware.JwtAuthMiddleware())
			roles.Use(middleware.RBACMiddleware(constants.RoleAdmin))
			{
				roles.POST("", roleHandler.CreateRole)
				roles.PUT("/:id", roleHandler.UpdateRole)
				roles.DELETE("/:id", roleHandler.DeleteRole)
				roles.GET("/:id", roleHandler.GetRoleByID)
				roles.GET("", roleHandler.GetPaginatedRoles)
			}
		}
	}
}
