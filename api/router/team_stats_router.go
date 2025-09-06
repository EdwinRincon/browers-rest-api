package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeTeamStatsRoutes(r *gin.Engine, teamStatsHandler *handler.TeamStatsHandler) {
	api := r.Group(constants.APIBasePath)
	auth := api.Group("")
	auth.Use(middleware.JwtAuthMiddleware())

	// Team Stats routes for authenticated users
	teamStats := auth.Group("/team-stats")
	{
		teamStats.GET("", teamStatsHandler.GetPaginatedTeamStats)
		teamStats.GET("/:id", teamStatsHandler.GetTeamStatsByID)
	}

	// Team-specific stats
	auth.GET("/teams/:id/stats", teamStatsHandler.GetTeamStatsByTeamID)

	// Season-specific stats
	auth.GET("/seasons/:id/team-stats", teamStatsHandler.GetTeamStatsBySeasonID)

	// Admin-only routes
	admin := api.Group("/admin/team-stats")
	admin.Use(middleware.JwtAuthMiddleware(), middleware.RBACMiddleware(constants.RoleAdmin))
	{
		admin.POST("", teamStatsHandler.CreateTeamStats)
		admin.PUT("/:id", teamStatsHandler.UpdateTeamStats)
		admin.DELETE("/:id", teamStatsHandler.DeleteTeamStats)
	}
}
