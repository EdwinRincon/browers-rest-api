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
		// Public match routes
		matches := api.Group("/matches")
		{
			matches.GET("", matchHandler.GetPaginatedMatches)             // GET /matches
			matches.GET("/:id", matchHandler.GetMatchByID)                // GET /matches/:id
			matches.GET("/:id/detail", matchHandler.GetDetailedMatchByID) // GET /matches/:id/detail
		}

		// Season-related match routes
		seasons := api.Group("/seasons")
		{
			seasons.GET("/:id/matches", matchHandler.GetMatchesBySeasonID) // GET /seasons/:id/matches
		}

		// Team-related match routes
		teams := api.Group("/teams")
		{
			teams.GET("/:id/matches", matchHandler.GetMatchesByTeamID)      // GET /teams/:id/matches
			teams.GET("/:id/next-match", matchHandler.GetNextMatchByTeamID) // GET /teams/:id/next-match
		}

		// Admin-only match routes
		admin := api.Group("/admin")
		admin.Use(middleware.JwtAuthMiddleware(), middleware.RBACMiddleware(constants.RoleAdmin))
		{
			adminMatches := admin.Group("/matches")
			{
				adminMatches.POST("", matchHandler.CreateMatch)       // POST /admin/matches
				adminMatches.PUT("/:id", matchHandler.UpdateMatch)    // PUT /admin/matches/:id
				adminMatches.DELETE("/:id", matchHandler.DeleteMatch) // DELETE /admin/matches/:id
			}
		}
	}
}
