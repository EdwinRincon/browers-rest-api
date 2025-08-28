package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeSeasonRoutes(r *gin.Engine, seasonHandler *handler.SeasonHandler) {
	api := r.Group(constants.APIBasePath)

	// Seasons endpoints (read-only, no authentication required)
	seasons := api.Group("/seasons")
	{
		seasons.GET("", seasonHandler.GetPaginatedSeasons)
		seasons.GET("/current", seasonHandler.GetCurrentSeason)
		seasons.GET("/:id", seasonHandler.GetSeasonByID)
	}

	// Admin routes (authenticated + role check)
	adminSeasons := api.Group("/admin/seasons")
	adminSeasons.Use(middleware.JwtAuthMiddleware(), middleware.RBACMiddleware(constants.RoleAdmin))
	{
		adminSeasons.POST("", seasonHandler.CreateSeason)
		adminSeasons.PUT("/:id", seasonHandler.UpdateSeason)
		adminSeasons.PUT("/:id/set-current", seasonHandler.SetCurrentSeason)
		adminSeasons.DELETE("/:id", seasonHandler.DeleteSeason)
	}
}
