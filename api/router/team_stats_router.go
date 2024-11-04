package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeTeamStatsRoutes(r *gin.Engine, teamStatsHandler *handler.TeamStatsHandler) {
	api := r.Group(constants.APIBasePath)
	{
		teamStats := api.Group("/team-stats")
		{
			teamStats.Use(middleware.JwtAuthMiddleware())

			teamStats.GET("", teamStatsHandler.ListTeamStats)
			teamStats.GET("/:id", teamStatsHandler.GetTeamStatsByID)

			teamStats.Use(middleware.RBACMiddleware(constants.RoleAdmin))
			{
				teamStats.POST("", teamStatsHandler.CreateTeamStats)
				teamStats.PUT("/:id", teamStatsHandler.UpdateTeamStats)
				teamStats.DELETE("/:id", teamStatsHandler.DeleteTeamStats)
			}
		}
	}
}
