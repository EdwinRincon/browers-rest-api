package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeMatchRoutes(r *gin.Engine, matchHandler *handler.MatchHandler) {
	r.Use(middleware.SecurityHeadersMiddleware())

	api := r.Group(constants.APIBasePath)
	{
		matches := api.Group("/matches")
		{
			matches.Use(middleware.JwtAuthMiddleware())

			matches.GET("", matchHandler.ListMatches)
			matches.GET("/:id", matchHandler.GetMatchByID)

			matches.Use(middleware.RBACMiddleware(constants.RoleAdmin, constants.RoleCoach))
			{
				matches.POST("", matchHandler.CreateMatch)
				matches.PUT("/:id", matchHandler.UpdateMatch)
				matches.DELETE("/:id", matchHandler.DeleteMatch)
			}
		}
	}
}
