package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeTeamRoutes(r *gin.Engine, teamHandler *handler.TeamHandler) {
	api := r.Group(constants.APIBasePath)

	// Authenticated user routes
	teams := api.Group("/teams")
	teams.Use(middleware.JwtAuthMiddleware())
	{
		teams.GET("", teamHandler.GetPaginatedTeams)
		teams.GET("/:id", teamHandler.GetTeamByID)
	}

	// Admin-only routes
	adminTeams := api.Group("/admin/teams")
	adminTeams.Use(middleware.JwtAuthMiddleware(), middleware.RBACMiddleware(constants.RoleAdmin))
	{
		adminTeams.POST("", teamHandler.CreateTeam)
		adminTeams.PUT("/:id", teamHandler.UpdateTeam)
		adminTeams.DELETE("/:id", teamHandler.DeleteTeam)
	}
}
