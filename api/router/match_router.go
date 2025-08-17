package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeMatchRoutes(r *gin.Engine, matchHandler *handler.MatchHandler) {

	api := r.Group(constants.APIBasePath)
	{
		matches := api.Group("/matches")
		{
			matches.GET("", matchHandler.ListMatches)

			matches.Use(middleware.JwtAuthMiddleware())
			matches.GET("/:id", matchHandler.GetMatchByID)

			matches.Use(middleware.RBACMiddleware(constants.RoleAdmin))
			{
				matches.POST("", matchHandler.CreateMatch)
				matches.PUT("/:id", matchHandler.UpdateMatch)
				matches.DELETE("/:id", matchHandler.DeleteMatch)
			}
		}
	}
}
