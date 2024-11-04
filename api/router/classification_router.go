package api

import (
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeClassificationRoutes(r *gin.Engine, classificationHandler *handler.ClassificationHandler) {
	r.Use(middleware.SecurityHeadersMiddleware())

	api := r.Group(constants.APIBasePath)
	{
		classifications := api.Group("/classifications")
		{
			classifications.Use(middleware.JwtAuthMiddleware())

			classifications.GET("", classificationHandler.ListClassifications)
			classifications.GET("/:id", classificationHandler.GetClassificationByID)
			classifications.GET("/season/:seasonID", classificationHandler.GetClassificationBySeason)

			classifications.Use(middleware.RBACMiddleware(constants.RoleAdmin))
			{
				classifications.POST("", classificationHandler.CreateClassification)
				classifications.PUT("/:id", classificationHandler.UpdateClassification)
				classifications.DELETE("/:id", classificationHandler.DeleteClassification)
			}
		}
	}
}
