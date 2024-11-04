package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeArticleRoutes(r *gin.Engine, articleHandler *handler.ArticleHandler) {
	r.Use(middleware.SecurityHeadersMiddleware())
	api := r.Group(constants.APIBasePath)
	{
		articles := api.Group("/articles")
		{
			articles.GET("", articleHandler.ListArticles)
			articles.GET("/:id", articleHandler.GetArticleByID)
			articles.Use(middleware.JwtAuthMiddleware())
			articles.Use(middleware.RBACMiddleware(constants.RoleAdmin))
			{
				articles.POST("", articleHandler.CreateArticle)
				articles.PUT("/:id", articleHandler.UpdateArticle)
				articles.DELETE("/:id", articleHandler.DeleteArticle)
			}
		}
	}
}
