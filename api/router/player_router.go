package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializePlayerRoutes(r *gin.Engine, playerHandler *handler.PlayerHandler) {

	api := r.Group(constants.APIBasePath)
	{
		players := api.Group("/players")
		{
			players.Use(middleware.JwtAuthMiddleware())

			players.GET("", playerHandler.GetAllPlayers)
			players.GET("/:id", playerHandler.GetPlayerByID)

			players.Use(middleware.RBACMiddleware(constants.RoleAdmin))
			{
				players.POST("", playerHandler.CreatePlayer)
				players.PUT("/:id", playerHandler.UpdatePlayer)
				players.DELETE("/:id", playerHandler.DeletePlayer)
			}
		}
	}
}
