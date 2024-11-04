package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeLineupRoutes(r *gin.Engine, lineupHandler *handler.LineupHandler) {
	r.Use(middleware.SecurityHeadersMiddleware())

	api := r.Group(constants.APIBasePath)
	{
		lineups := api.Group("/lineups")
		{
			lineups.Use(middleware.JwtAuthMiddleware())

			lineups.GET("", lineupHandler.ListLineups)
			lineups.GET("/:id", lineupHandler.GetLineupByID)
			lineups.GET("/match/:matchID", lineupHandler.GetLineupsByMatch)

			lineups.Use(middleware.RBACMiddleware(constants.RoleAdmin))
			{
				lineups.POST("", lineupHandler.CreateLineup)
				lineups.PUT("/:id", lineupHandler.UpdateLineup)
				lineups.DELETE("/:id", lineupHandler.DeleteLineup)
			}
		}
	}
}
