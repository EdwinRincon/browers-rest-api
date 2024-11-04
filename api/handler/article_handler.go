package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
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
func (h ArticleHandler) CreateArticle(c *gin.Context) {
	var article model.Articles
	if err := c.ShouldBindJSON(&article); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}
	ctx := c.Request.Context()
	err := h.ArticleService.CreateArticle(ctx, &article)
	if err != nil {
		if errors.Is(err, repository.ErrArticleNotFound) {
			helper.HandleError(c, helper.NewAppError(http.StatusNotFound, constants.ErrArticleNotFound.Error(), ""), false)
			return
		}
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to retrieve article", err.Error()), false)
		return
	}
	helper.HandleSuccess(c, http.StatusCreated, article, "Article created successfully")
}
func (h ArticleHandler) GetArticleByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidArticleID, err.Error()), true)
		return
	}
	ctx := c.Request.Context()
	article, err := h.ArticleService.GetArticleByID(ctx, id)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusNotFound, "Article not found", err.Error()), false)
		return
	}
	helper.HandleSuccess(c, http.StatusOK, article, "Article retrieved successfully")
}

func (h ArticleHandler) ListArticles(c *gin.Context) {
	page, err := strconv.ParseUint(c.DefaultQuery("page", "1"), 10, 64)
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
	articles, err := h.ArticleService.ListArticles(ctx, page, pageSize)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to list articles", err.Error()), true)
		return
	}
	helper.HandleSuccess(c, http.StatusOK, articles, "Articles listed successfully")
}
func (h ArticleHandler) UpdateArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidArticleID, err.Error()), true)
		return
	}
	var article model.Articles
	if err := c.ShouldBindJSON(&article); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}
	article.ID = id
	ctx := c.Request.Context()
	err = h.ArticleService.UpdateArticle(ctx, &article)
	if err != nil {
		if errors.Is(err, repository.ErrArticleNotFound) {
			helper.HandleError(c, helper.NewAppError(http.StatusNotFound, constants.ErrArticleNotFound.Error(), ""), true)
			return
		}
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to update article", err.Error()), true)
		return
	}
	helper.HandleSuccess(c, http.StatusOK, article, "Article updated successfully")
}

func (h ArticleHandler) DeleteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidArticleID, err.Error()), true)
		return
	}
	ctx := c.Request.Context()
	err = h.ArticleService.DeleteArticle(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrArticleNotFound) {
			helper.HandleError(c, helper.NewAppError(http.StatusNotFound, constants.ErrArticleNotFound.Error(), ""), true)
			return
		}
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to delete article", err.Error()), true)
		return
	}
	helper.HandleSuccess(c, http.StatusOK, nil, "Article deleted successfully")
}
