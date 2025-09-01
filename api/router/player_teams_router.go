package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializePlayerTeamRoutes(r *gin.Engine, playerTeamHandler *handler.PlayerTeamHandler) {
	api := r.Group(constants.APIBasePath)

	// All authenticated routes
	authenticatedRoutes := api.Group("")
	authenticatedRoutes.Use(middleware.JwtAuthMiddleware())
	{
		// player-team operations
		playerTeams := authenticatedRoutes.Group("/player-teams")
		{
			playerTeams.GET("", playerTeamHandler.GetPaginatedPlayerTeams)
			playerTeams.GET("/:id", playerTeamHandler.GetPlayerTeamByID)
		}

		// Player-team relationship views by different entities
		authenticatedRoutes.GET("/players/:id/teams", playerTeamHandler.GetPlayerTeamsByPlayerID)
		authenticatedRoutes.GET("/teams/:id/players", playerTeamHandler.GetPlayerTeamsByTeamID)
		authenticatedRoutes.GET("/seasons/:id/player-teams", playerTeamHandler.GetPlayerTeamsBySeasonID)
	}

	// Admin-only
	adminRoutes := api.Group("/admin")
	adminRoutes.Use(middleware.JwtAuthMiddleware(), middleware.RBACMiddleware(constants.RoleAdmin))
	{
		admin := adminRoutes.Group("/player-teams")
		{
			admin.POST("", playerTeamHandler.CreatePlayerTeam)
			admin.PUT("/:id", playerTeamHandler.UpdatePlayerTeam)
			admin.DELETE("/:id", playerTeamHandler.DeletePlayerTeam)
		}
	}
}
