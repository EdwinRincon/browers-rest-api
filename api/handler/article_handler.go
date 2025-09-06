package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/mapper"
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

// CreateArticle godoc
// @Summary      Create a new article
// @Description  Creates a new article with the provided data
// @Tags         articles
// @ID           createArticle
// @Accept       json
// @Produce      json
// @Param        article  body      dto.CreateArticleRequest  true  "Article data"
// @Success      201   {object}  dto.ArticleShort  "Article created successfully"
// @Failure      400   {object}  helper.AppError "Invalid input"
// @Failure      404   {object}  helper.AppError "Season not found"
// @Failure      500   {object}  helper.AppError "Internal server error"
// @Router       /admin/articles [post]
// @Security     ApiKeyAuth
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var createRequest dto.CreateArticleRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		helper.WriteErrorResponse(c, helper.BuildValidationErrorFromBinding(err, "body", "Invalid article data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Map the request to an Article model
	article := mapper.ToArticle(&createRequest)

	createdArticle, err := h.ArticleService.CreateArticle(ctx, article)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrSeasonNotFound):
			helper.WriteErrorResponse(c, helper.NewBadRequestError("season_id", "Season not found"))
			return
		default:
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
			return
		}
	}

	helper.WriteSuccessResponse(c, http.StatusCreated, createdArticle, "Article created successfully")
}

// GetArticleByID godoc
// @Summary      Get an article by ID
// @Description  Retrieves an article by its ID
// @Tags         articles
// @ID           getArticleByID
// @Param        id  path      int  true  "Article ID"
// @Success      200 {object}  dto.ArticleResponse "Article retrieved successfully"
// @Failure      400 {object}  helper.AppError "Invalid ID format"
// @Failure      404 {object}  helper.AppError "Article not found"
// @Failure      500 {object}  helper.AppError "Internal server error"
// @Router       /articles/{id} [get]
func (h *ArticleHandler) GetArticleByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", constants.MsgInvalidID))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	article, err := h.ArticleService.GetArticleByID(ctx, id)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("article"))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	articleResponse := mapper.ToArticleResponse(article)
	helper.WriteSuccessResponse(c, http.StatusOK, articleResponse, "Article retrieved successfully")
}

// GetPaginatedArticles godoc
// @Summary      Get paginated articles
// @Description  Retrieves a paginated list of articles with sorting and ordering
// @Tags         articles
// @ID           getPaginatedArticles
// @Param        page      query     int     false  "Page number (0-based)"
// @Param        pageSize  query     int     false  "Number of items per page (default 10)"
// @Param        sort      query     string  false  "Sort field (e.g., title, date, created_at)"
// @Param        order     query     string  false  "Sort order (asc/desc)"
// @Success      200       {object}  map[string]interface{} "Articles retrieved successfully"
// @Failure      400       {object}  helper.AppError "Invalid input"
// @Failure      500       {object}  helper.AppError "Internal server error"
// @Router       /articles [get]
func (h *ArticleHandler) GetPaginatedArticles(c *gin.Context) {
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 0 {
		page = 0
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 10
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	// Validate sort field
	if err := helper.ValidateSort(model.Article{}, sort); err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("sort", err.Error()))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	articles, total, err := h.ArticleService.GetPaginatedArticles(ctx, sort, order, page, pageSize)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	response := helper.PaginatedResponse{
		Items:      mapper.ToArticleResponseList(articles),
		TotalCount: total,
	}

	helper.WriteSuccessResponse(c, http.StatusOK, response, "Articles retrieved successfully")
}

// GetArticlesBySeasonID godoc
// @Summary      Get articles by season ID
// @Description  Retrieves a paginated list of articles for a specific season
// @Tags         articles
// @ID           getArticlesBySeasonID
// @Param        id        path      int     true   "Season ID"
// @Param        page      query     int     false  "Page number (0-based)"
// @Param        pageSize  query     int     false  "Number of items per page (default 10)"
// @Param        sort      query     string  false  "Sort field (e.g., title, date, created_at)"
// @Param        order     query     string  false  "Sort order (asc/desc)"
// @Success      200       {object}  map[string]interface{} "Articles retrieved successfully"
// @Failure      400       {object}  helper.AppError "Invalid input"
// @Failure      404       {object}  helper.AppError "Season not found"
// @Failure      500       {object}  helper.AppError "Internal server error"
// @Router       /seasons/{id}/articles [get]
func (h *ArticleHandler) GetArticlesBySeasonID(c *gin.Context) {
	seasonIDStr := c.Param("id")
	seasonID, err := strconv.ParseUint(seasonIDStr, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid season ID format"))
		return
	}

	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 0 {
		page = 0
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 10
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	// Validate sort field
	if err := helper.ValidateSort(model.Article{}, sort); err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("sort", err.Error()))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	articles, total, err := h.ArticleService.GetArticlesBySeasonID(ctx, seasonID, sort, order, page, pageSize)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrSeasonNotFound):
			helper.WriteErrorResponse(c, helper.NewNotFoundError("season"))
			return
		default:
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
			return
		}
	}

	response := helper.PaginatedResponse{
		Items:      mapper.ToArticleResponseList(articles),
		TotalCount: total,
	}

	helper.WriteSuccessResponse(c, http.StatusOK, response, "Articles retrieved successfully")
}

// UpdateArticle godoc
// @Summary      Update an existing article
// @Description  Updates an existing article's information by ID
// @Tags         articles
// @ID           updateArticle
// @Accept       json
// @Produce      json
// @Param        id       path      int                      true  "Article ID"
// @Param        article  body      dto.UpdateArticleRequest true  "Updated article data"
// @Success      200      {object}  dto.ArticleShort  "Article updated successfully"
// @Failure      400      {object}  helper.AppError "Invalid input or ID format"
// @Failure      404      {object}  helper.AppError "Article or Season not found"
// @Failure      500      {object}  helper.AppError "Internal server error"
// @Router       /admin/articles/{id} [put]
// @Security     ApiKeyAuth
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", constants.MsgInvalidID))
		return
	}

	var articleUpdateDTO dto.UpdateArticleRequest
	if err := c.ShouldBindJSON(&articleUpdateDTO); err != nil {
		helper.WriteErrorResponse(c, helper.BuildValidationErrorFromBinding(err, "body", "Invalid article data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	updatedArticle, err := h.ArticleService.UpdateArticle(ctx, &articleUpdateDTO, id)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrRecordNotFound):
			helper.WriteErrorResponse(c, helper.NewNotFoundError("article"))
			return
		case errors.Is(err, constants.ErrSeasonNotFound):
			helper.WriteErrorResponse(c, helper.NewBadRequestError("season_id", "Season not found"))
			return
		default:
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
			return
		}
	}

	response := mapper.ToArticleShort(updatedArticle)
	helper.WriteSuccessResponse(c, http.StatusOK, response, "Article updated successfully")
}

// DeleteArticle godoc
// @Summary      Delete an article
// @Description  Deletes an article by its ID
// @Tags         articles
// @ID           deleteArticle
// @Param        id   path      int  true  "Article ID"
// @Success      204 "No Content"
// @Failure      400  {object}  helper.AppError "Invalid ID format"
// @Router       /admin/articles/{id} [delete]
// @Security     ApiKeyAuth
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", constants.MsgInvalidID))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.ArticleService.DeleteArticle(ctx, id)
	if err != nil && !errors.Is(err, constants.ErrRecordNotFound) {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	c.Status(http.StatusNoContent)
}
