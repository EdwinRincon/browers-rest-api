package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/EdwinRincon/browersfc-api/internal/domain/service"
	"github.com/gin-gonic/gin"
)

func InitializePlayerRoutes(r *gin.Engine, playerHandler *handler.PlayerHandler, authService *service.AuthenticationDomainService) {
	api := r.Group(constants.APIBasePath)
	{
		// Public/Authenticated player routes
		players := api.Group("/players")
		players.Use(middleware.JwtAuthMiddleware(authService))
		{
			players.GET("", playerHandler.GetPaginatedPlayers)
			players.GET("/nickname/:nickname", playerHandler.GetPlayerByNickName) // placed before :id
			players.GET("/:id", playerHandler.GetPlayerByID)
		}

		// Admin routes
		adminPlayers := api.Group("/admin/players")
		adminPlayers.Use(middleware.JwtAuthMiddleware(authService), middleware.RBACMiddleware(constants.RoleAdmin))
		{
			adminPlayers.POST("", playerHandler.CreatePlayer)
			adminPlayers.PUT("/:id", playerHandler.UpdatePlayer)
			adminPlayers.DELETE("/:id", playerHandler.DeletePlayer)
		}
	}
}
