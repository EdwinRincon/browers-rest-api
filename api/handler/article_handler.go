package handler

import (
	"net/http"
	"strconv"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/service"
	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	ArticleService service.ArticleService
}

func NewArticleHandler(articleService service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		ArticleService: articleService,
	}
}

func (h *ArticleHandler) GetArticleByID(c *gin.Context) {
	articleID := c.Param("id")
	id, err := strconv.ParseUint(articleID, 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	article, err := h.ArticleService.GetArticleByID(ctx, id)
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, article, "Article found successfully")
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var article model.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		helper.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.ArticleService.CreateArticle(ctx, &article)
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusCreated, article, "Article created successfully")
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	articleID := c.Param("id")
	id, err := strconv.ParseUint(articleID, 10, 32)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	var article model.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		helper.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	article.ID = uint(id)
	err = h.ArticleService.UpdateArticle(ctx, &article)
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, article, "Article updated successfully")
}

func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	articleID := c.Param("id")
	id, err := strconv.ParseUint(articleID, 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err = h.ArticleService.DeleteArticle(ctx, id)
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, nil, "Article deleted successfully")
}

func (h *ArticleHandler) GetAllArticles(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1") // Default to page 1 instead of 0

	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid page number", err.Error()), true)
		return
	}
	pageSize, err := strconv.ParseUint(c.DefaultQuery("pageSize", "10"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid page size", err.Error()), true)
		return
	}
	ctx := c.Request.Context()
	articles, err := h.ArticleService.GetAllArticles(ctx, page, pageSize)
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, articles, "Article retrieved successfully")
}
