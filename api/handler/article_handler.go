package handler

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"time"

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

// GetArticleByID godoc
// @Summary      Get an article by ID
// @Description  Returns the details of an article by its ID
// @Tags         articles
// @ID           getArticleByID
// @Param        id   path      int  true  "Article ID"
// @Success      200  {object}  model.Article "Success"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Article not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /articles/{id} [get]
// @Security     ApiKeyAuth
func (h *ArticleHandler) GetArticleByID(c *gin.Context) {
	articleID := c.Param("id")
	id, err := strconv.ParseUint(articleID, 10, 64)
	if err != nil || id > uint64(math.MaxUint32) {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid or too large article ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	article, err := h.ArticleService.GetArticleByID(ctx, id)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	helper.WriteSuccessResponse(c, http.StatusOK, article, "Article found successfully")
}

// CreateArticle godoc
// @Summary      Create a new article
// @Description  Creates a new article with the provided data
// @Tags         articles
// @ID           createArticle
// @Accept       json
// @Produce      json
// @Param        article  body      model.Article  true  "Article data"
// @Success      201      {object}  model.Article "Created"
// @Failure      400      {object}  helper.AppError "Invalid input"
// @Failure      409      {object}  helper.AppError "Conflict (e.g., article already exists)"
// @Failure      500      {object}  helper.AppError "Internal server error"
// @Router       /articles [post]
// @Security     ApiKeyAuth
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var article model.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("body", "Invalid article data"))
		return
	}
	article.ID = 0

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	err := h.ArticleService.CreateArticle(ctx, &article)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	helper.WriteSuccessResponse(c, http.StatusCreated, article, "Article created successfully")
}

// UpdateArticle godoc
// @Summary      Update an existing article
// @Description  Updates the details of an existing article by ID
// @Tags         articles
// @ID           updateArticle
// @Accept       json
// @Produce      json
// @Param        id       path      int           true  "Article ID"
// @Param        article  body      model.Article true  "Updated article data"
// @Success      200      {object}  model.Article "Updated"
// @Failure      400      {object}  helper.AppError "Invalid input"
// @Failure      404      {object}  helper.AppError "Article not found"
// @Failure      500      {object}  helper.AppError "Internal server error"
// @Router       /articles/{id} [put]
// @Security     ApiKeyAuth
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	articleID := c.Param("id")
	id, err := strconv.ParseUint(articleID, 10, 64)
	if err != nil || id > uint64(math.MaxUint32) {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid or too large article ID"))
		return
	}

	var article model.Article
	if err = c.ShouldBindJSON(&article); err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("body", "Invalid article data"))
		return
	}

	if article.ID != 0 && article.ID != id {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Mismatched article ID in body and URL"))
		return
	}

	article.ID = id

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	err = h.ArticleService.UpdateArticle(ctx, &article)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	helper.WriteSuccessResponse(c, http.StatusOK, article, "Article updated successfully")
}

// DeleteArticle godoc
// @Summary      Delete an article
// @Description  Deletes an article by its ID
// @Tags         articles
// @ID           deleteArticle
// @Param        id   path      int  true  "Article ID"
// @Success      204 "No Content"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Article not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /articles/{id} [delete]
// @Security     ApiKeyAuth
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	articleID := c.Param("id")
	id, err := strconv.ParseUint(articleID, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid article ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	err = h.ArticleService.DeleteArticle(ctx, id)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAllArticles godoc
// @Summary      List articles
// @Description  Retrieves a paginated list of articles
// @Tags         articles
// @ID           listArticles
// @Param        page      query     int  false  "Page number"
// @Param        pageSize  query     int  false  "Page size"
// @Success      200       {array}   model.Article "Articles listed successfully"
// @Failure      400       {object}  helper.AppError "Invalid input"
// @Failure      500       {object}  helper.AppError "Internal server error"
// @Router       /articles [get]
// @Security     ApiKeyAuth
func (h *ArticleHandler) GetAllArticles(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("page", "Invalid page number"))
		return
	}
	pageSize, err := strconv.ParseUint(pageSizeStr, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("pageSize", "Invalid page size"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	articles, err := h.ArticleService.GetAllArticles(ctx, page, pageSize)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	helper.WriteSuccessResponse(c, http.StatusOK, articles, "Articles retrieved successfully")
}
