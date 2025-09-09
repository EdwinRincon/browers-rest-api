package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/EdwinRincon/browersfc-api/internal/domain/service"
	"github.com/gin-gonic/gin"
)

func InitializeLineupRoutes(r *gin.Engine, lineupHandler *handler.LineupHandler, authService *service.AuthenticationDomainService) {
	api := r.Group(constants.APIBasePath)

	authRequired := middleware.JwtAuthMiddleware(authService)

	lineups := api.Group("/lineups", authRequired)
	{
		lineups.GET("", lineupHandler.GetPaginatedLineups)
		lineups.GET("/:id", lineupHandler.GetLineupByID)
	}

	// Specific lineup queries by match type
	lineupsByMatch := api.Group("/lineups/match", authRequired)
	{
		lineupsByMatch.GET("/:id/starting", lineupHandler.GetStartingLineupsByMatchID)
		lineupsByMatch.GET("/:id/substitutes", lineupHandler.GetSubstitutesLineupsByMatchID)
	}

	matchLineups := api.Group("/matches/:id/lineups", authRequired)
	{
		matchLineups.GET("", lineupHandler.GetLineupsByMatchID)
	}

	playerLineups := api.Group("/players/:id/lineups", authRequired)
	{
		playerLineups.GET("", lineupHandler.GetLineupsByPlayerID)
	}

	// --- Admin-only lineup management ---
	adminLineups := api.Group("/admin/lineups", authRequired, middleware.RBACMiddleware(constants.RoleAdmin))
	{
		adminLineups.POST("", lineupHandler.CreateLineup)
		adminLineups.PUT("/:id", lineupHandler.UpdateLineup)
		adminLineups.DELETE("/:id", lineupHandler.DeleteLineup)
	}
}
