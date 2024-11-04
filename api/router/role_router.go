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
			// Rutas accesibles por todos los usuarios autenticados
			roles.GET("/:id", roleHandler.GetRoleByID)
			roles.GET("", roleHandler.GetAllRoles)

			// Rutas accesibles solo por administradores
			roles.Use(middleware.RBACMiddleware(constants.RoleAdmin)) // Ensure this is applied before sensitive routes
			{
				roles.POST("", roleHandler.CreateRole)
				roles.PUT("/:id", roleHandler.UpdateRole)
				roles.DELETE("/:id", roleHandler.DeleteRole) // This should be protected
			}

		}
	}
}
