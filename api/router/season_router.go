package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeSeasonRoutes(r *gin.Engine, seasonHandler *handler.SeasonHandler) {
	r.Use(middleware.SecurityHeadersMiddleware())

	api := r.Group(constants.APIBasePath)
	{
		seasons := api.Group("/seasons")
		{
			seasons.Use(middleware.JwtAuthMiddleware())

			seasons.GET("", seasonHandler.ListSeasons)
			seasons.GET("/:id", seasonHandler.GetSeasonByID)

			seasons.Use(middleware.RBACMiddleware(constants.RoleAdmin))
			{
				seasons.POST("", seasonHandler.CreateSeason)
				seasons.PUT("/:id", seasonHandler.UpdateSeason)
				seasons.DELETE("/:id", seasonHandler.DeleteSeason)
			}
		}
	}
}
