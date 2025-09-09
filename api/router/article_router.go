package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/EdwinRincon/browersfc-api/internal/domain/service"
	"github.com/gin-gonic/gin"
)

func InitializeArticleRoutes(r *gin.Engine, articleHandler *handler.ArticleHandler, authService *service.AuthenticationDomainService) {
	api := r.Group(constants.APIBasePath)
	{

		articles := api.Group("/articles")
		{
			articles.GET("", articleHandler.GetPaginatedArticles)
			articles.GET("/:id", articleHandler.GetArticleByID)
		}

		seasons := api.Group("/seasons")
		{
			seasons.GET("/:id/articles", articleHandler.GetArticlesBySeasonID)
		}

		// Articles routes requiring authentication and role-based access control
		protected := api.Group("/admin/articles")
		protected.Use(middleware.JwtAuthMiddleware(authService), middleware.RBACMiddleware(constants.RoleAdmin))
		{
			protected.POST("", articleHandler.CreateArticle)
			protected.PUT("/:id", articleHandler.UpdateArticle)
			protected.DELETE("/:id", articleHandler.DeleteArticle)
		}
	}
}
