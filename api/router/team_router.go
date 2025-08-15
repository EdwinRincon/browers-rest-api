package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeTeamRoutes(r *gin.Engine, teamHandler *handler.TeamHandler) {
	r.Use(middleware.SecurityHeadersMiddleware())

	api := r.Group(constants.APIBasePath)
	{
		teams := api.Group("/teams")
		{
			teams.Use(middleware.JwtAuthMiddleware())

			teams.GET("", teamHandler.GetPaginatedTeams)
			teams.GET("/:id", teamHandler.GetTeamByID)

			teams.Use(middleware.RBACMiddleware(constants.RoleAdmin))
			{
				teams.POST("", teamHandler.CreateTeam)
				teams.PUT("/:id", teamHandler.UpdateTeam)
				teams.DELETE("/:id", teamHandler.DeleteTeam)
			}
		}
	}
}
