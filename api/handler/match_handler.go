package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/service"
	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MatchHandler struct {
	MatchService service.MatchService
}

func NewMatchHandler(matchService service.MatchService) *MatchHandler {
	return &MatchHandler{
		MatchService: matchService,
	}
}

// CreateMatch godoc
// @Summary      Create a new match
// @Description  Creates a new match with the provided data
// @Tags         matches
// @ID           createMatch
// @Accept       json
// @Produce      json
// @Param        match  body      model.Match  true  "Match data"
// @Success      201    {object}  model.Match "Created"
// @Failure      400    {object}  helper.AppError "Invalid input"
// @Failure      409    {object}  helper.AppError "Conflict (e.g., match already exists)"
// @Failure      500    {object}  helper.AppError "Internal server error"
// @Router       /matches [post]
// @Security     ApiKeyAuth
func (h *MatchHandler) CreateMatch(c *gin.Context) {
	var match model.Match
	if err := c.ShouldBindJSON(&match); err != nil {
		helper.RespondWithError(c, helper.BadRequest("body", "Invalid match data"))
		return
	}

	match.ID = 0
	ctx := c.Request.Context()

	if err := h.MatchService.CreateMatch(ctx, &match); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			helper.RespondWithError(c, helper.Conflict("match", "A match with these details already exists"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusCreated, match, "Match created successfully")
}

// GetMatchByID godoc
// @Summary      Get a match by ID
// @Description  Returns the details of a match by its ID
// @Tags         matches
// @ID           getMatchByID
// @Param        id   path      int  true  "Match ID"
// @Success      200  {object}  model.Match "Success"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Match not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /matches/{id} [get]
// @Security     ApiKeyAuth
func (h *MatchHandler) GetMatchByID(c *gin.Context) {
	matchID := c.Param("id")
	id, err := strconv.ParseUint(matchID, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid match ID"))
		return
	}

	ctx := c.Request.Context()
	match, err := h.MatchService.GetMatchByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("match"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusOK, match, "Match found successfully")
}

// ListMatches godoc
// @Summary      List matches
// @Description  Retrieves a paginated list of matches
// @Tags         matches
// @ID           listMatches
// @Param        page      query     int  false  "Page number"
// @Param        pageSize  query     int  false  "Page size"
// @Success      200       {array}   model.Match "Matches listed successfully"
// @Failure      400       {object}  helper.AppError "Invalid input"
// @Failure      500       {object}  helper.AppError "Internal server error"
// @Router       /matches [get]
// @Security     ApiKeyAuth
func (h *MatchHandler) ListMatches(c *gin.Context) {
	page, err := strconv.ParseUint(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("page", "Invalid page number"))
		return
	}

	pageSize, err := strconv.ParseUint(c.DefaultQuery("pageSize", "10"), 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("pageSize", "Invalid page size"))
		return
	}

	ctx := c.Request.Context()
	matches, err := h.MatchService.ListMatches(ctx, page, pageSize)
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	helper.HandleSuccess(c, http.StatusOK, matches, "Match listed successfully")
}

// UpdateMatch godoc
// @Summary      Update an existing match
// @Description  Updates the details of an existing match by ID
// @Tags         matches
// @ID           updateMatch
// @Accept       json
// @Produce      json
// @Param        id     path      int        true  "Match ID"
// @Param        match  body      model.Match true  "Updated match data"
// @Success      200    {object}  model.Match "Updated"
// @Failure      400    {object}  helper.AppError "Invalid input"
// @Failure      404    {object}  helper.AppError "Match not found"
// @Failure      500    {object}  helper.AppError "Internal server error"
// @Router       /matches/{id} [put]
// @Security     ApiKeyAuth
func (h *MatchHandler) UpdateMatch(c *gin.Context) {
	matchID := c.Param("id")
	id, err := strconv.ParseUint(matchID, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid match ID"))
		return
	}

	var match model.Match
	if err := c.ShouldBindJSON(&match); err != nil {
		helper.RespondWithError(c, helper.BadRequest("body", "Invalid match data"))
		return
	}

	ctx := c.Request.Context()
	match.ID = id
	err = h.MatchService.UpdateMatch(ctx, &match)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("match"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusOK, match, "Match updated successfully")
}

// DeleteMatch godoc
// @Summary      Delete a match
// @Description  Deletes a match by its ID
// @Tags         matches
// @ID           deleteMatch
// @Param        id   path      int  true  "Match ID"
// @Success      204 nil "No Content"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Match not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /matches/{id} [delete]
// @Security     ApiKeyAuth
func (h *MatchHandler) DeleteMatch(c *gin.Context) {
	matchID := c.Param("id")
	id, err := strconv.ParseUint(matchID, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid match ID"))
		return
	}

	ctx := c.Request.Context()
	err = h.MatchService.DeleteMatch(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("match"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	c.Status(http.StatusNoContent)
}
