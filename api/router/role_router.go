package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/EdwinRincon/browersfc-api/internal/domain/service"
	"github.com/gin-gonic/gin"
)

func InitializeRoleRoutes(r *gin.Engine, roleHandler *handler.RoleHandler, authService *service.AuthenticationDomainService) {
	api := r.Group(constants.APIBasePath)

	// Admin-only role management
	adminRoles := api.Group("/admin/roles")
	adminRoles.Use(middleware.JwtAuthMiddleware(authService), middleware.RBACMiddleware(constants.RoleAdmin))
	// Read
	adminRoles.GET("", roleHandler.GetPaginatedRoles)
	adminRoles.GET("/:id", roleHandler.GetRoleByID)
	// Write
	adminRoles.POST("", roleHandler.CreateRole)
	adminRoles.PUT("/:id", roleHandler.UpdateRole)
	adminRoles.DELETE("/:id", roleHandler.DeleteRole)
}
