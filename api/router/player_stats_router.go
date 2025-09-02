package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializePlayerStatsRoutes(r *gin.Engine, playerStatsHandler *handler.PlayerStatsHandler) {
	api := r.Group(constants.APIBasePath)
	{
		// Public routes
		playerStats := api.Group("/player-stats")
		{
			playerStats.GET("", playerStatsHandler.GetPaginatedPlayerStats)
			playerStats.GET("/:id", playerStatsHandler.GetPlayerStatByID)
		}

		// Player-specific stats
		players := api.Group("/players")
		{
			players.GET("/:id/stats", playerStatsHandler.GetPlayerStatsByPlayerID)
		}

		// Match-specific stats
		matches := api.Group("/matches")
		{
			matches.GET("/:id/stats", playerStatsHandler.GetPlayerStatsByMatchID)
		}

		// Season-specific stats
		seasons := api.Group("/seasons")
		{
			seasons.GET("/:id/stats", playerStatsHandler.GetPlayerStatsBySeasonID)
		}

		// Admin routes
		adminPlayerStats := api.Group("/admin/player-stats")
		adminPlayerStats.Use(middleware.JwtAuthMiddleware(), middleware.RBACMiddleware(constants.RoleAdmin))
		{
			adminPlayerStats.POST("", playerStatsHandler.CreatePlayerStat)
			adminPlayerStats.PUT("/:id", playerStatsHandler.UpdatePlayerStat)
			adminPlayerStats.DELETE("/:id", playerStatsHandler.DeletePlayerStat)
		}
	}
}
