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
// @Param        match  body      dto.CreateMatchRequest  true  "Match data"
// @Success      201    {object}  dto.MatchResponse "Created match"
// @Failure      400    {object}  helper.AppError "Invalid input"
// @Failure      409    {object}  helper.AppError "Conflict"
// @Failure      500    {object}  helper.AppError "Internal server error"
// @Router       /admin/matches [post]
// @Security     ApiKeyAuth
func (h *MatchHandler) CreateMatch(c *gin.Context) {
	var createRequest dto.CreateMatchRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		helper.RespondWithError(c, helper.ProcessValidationError(err, "body", "Invalid match data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	createdMatch, err := h.MatchService.CreateMatch(ctx, &createRequest)
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	matchResponse := mapper.ToMatchShort(createdMatch)
	helper.HandleSuccess(c, http.StatusCreated, matchResponse, "Match created successfully")
}

// GetMatchByID godoc
// @Summary      Get a match by ID
// @Description  Returns the details of a match by its ID
// @Tags         matches
// @ID           getMatchByID
// @Param        id   path      int  true  "Match ID"
// @Success      200  {object}  dto.MatchResponse "Match found"
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

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	match, err := h.MatchService.GetMatchByID(ctx, id)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("match"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	matchResponse := mapper.ToMatchResponse(match)
	helper.HandleSuccess(c, http.StatusOK, matchResponse, "Match found successfully")
}

// GetDetailedMatchByID godoc
// @Summary      Get detailed match information
// @Description  Returns the complete details of a match including lineups and stats
// @Tags         matches
// @ID           getDetailedMatchByID
// @Param        id   path      int  true  "Match ID"
// @Success      200  {object}  dto.MatchDetailResponse "Match details found"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Match not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /matches/{id}/detail [get]
// @Security     ApiKeyAuth
func (h *MatchHandler) GetDetailedMatchByID(c *gin.Context) {
	matchID := c.Param("id")
	id, err := strconv.ParseUint(matchID, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid match ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	match, err := h.MatchService.GetDetailedMatchByID(ctx, id)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("match"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	detailedResponse := mapper.ToMatchDetailResponse(match)
	helper.HandleSuccess(c, http.StatusOK, detailedResponse, "Match details found successfully")
}

// GetPaginatedMatches godoc
// @Summary      Get paginated matches
// @Description  Retrieves a paginated list of matches with sorting and ordering
// @Tags         matches
// @ID           getPaginatedMatches
// @Param        page      query     int     false  "Page number (0-based)"
// @Param        pageSize  query     int     false  "Number of items per page (default 10)"
// @Param        sort      query     string  false  "Sort field (e.g., date, status)"
// @Param        order     query     string  false  "Sort order (asc/desc)"
// @Success      200       {object}  map[string]interface{} "Matches retrieved successfully"
// @Failure      400       {object}  helper.AppError "Invalid input"
// @Failure      500       {object}  helper.AppError "Internal server error"
// @Router       /matches [get]
func (h *MatchHandler) GetPaginatedMatches(c *gin.Context) {
	sort := c.DefaultQuery("sort", "date")
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
		order = "desc"
	}

	// Validate sort field
	if err := helper.ValidateSort(model.Match{}, sort); err != nil {
		helper.RespondWithError(c, helper.BadRequest("sort", err.Error()))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	matches, total, err := h.MatchService.GetPaginatedMatches(ctx, sort, order, page, pageSize)
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	response := helper.PaginatedResponse{
		Items:      mapper.ToMatchResponseList(matches),
		TotalCount: total,
	}

	helper.HandleSuccess(c, http.StatusOK, response, "Matches retrieved successfully")
}

// GetMatchesBySeasonID godoc
// @Summary      Get matches by season
// @Description  Retrieves matches for a specific season with pagination
// @Tags         matches
// @ID           getMatchesBySeasonID
// @Param        id        path      int     true   "Season ID"
// @Param        page      query     int     false  "Page number (0-based)"
// @Param        pageSize  query     int     false  "Page size (default 10)"
// @Param        sort      query     string  false  "Sort field"
// @Param        order     query     string  false  "Sort order (asc/desc)"
// @Success      200       {object}  map[string]interface{} "Matches retrieved successfully"
// @Failure      400       {object}  helper.AppError "Invalid input"
// @Failure      500       {object}  helper.AppError "Internal server error"
// @Router       /seasons/{id}/matches [get]
func (h *MatchHandler) GetMatchesBySeasonID(c *gin.Context) {
	seasonIDStr := c.Param("id")
	seasonID, err := strconv.ParseUint(seasonIDStr, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid season ID"))
		return
	}

	sort := c.DefaultQuery("sort", "date")
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
		order = "desc"
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	matches, total, err := h.MatchService.GetMatchesBySeasonID(ctx, seasonID, sort, order, page, pageSize)
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	response := helper.PaginatedResponse{
		Items:      mapper.ToMatchResponseList(matches),
		TotalCount: total,
	}

	helper.HandleSuccess(c, http.StatusOK, response, "Season matches retrieved successfully")
}

// GetMatchesByTeamID godoc
// @Summary      Get matches by team
// @Description  Retrieves matches where a specific team plays (home or away)
// @Tags         matches
// @ID           getMatchesByTeamID
// @Param        id        path      int     true   "Team ID"
// @Param        page      query     int     false  "Page number (0-based)"
// @Param        pageSize  query     int     false  "Page size (default 10)"
// @Param        sort      query     string  false  "Sort field"
// @Param        order     query     string  false  "Sort order (asc/desc)"
// @Success      200       {object}  map[string]interface{} "Matches retrieved successfully"
// @Failure      400       {object}  helper.AppError "Invalid input"
// @Failure      500       {object}  helper.AppError "Internal server error"
// @Router       /teams/{id}/matches [get]
func (h *MatchHandler) GetMatchesByTeamID(c *gin.Context) {
	teamIDStr := c.Param("id")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid team ID"))
		return
	}

	sort := c.DefaultQuery("sort", "date")
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
		order = "desc"
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	matches, total, err := h.MatchService.GetMatchesByTeamID(ctx, teamID, sort, order, page, pageSize)
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	response := helper.PaginatedResponse{
		Items:      mapper.ToMatchResponseList(matches),
		TotalCount: total,
	}

	helper.HandleSuccess(c, http.StatusOK, response, "Team matches retrieved successfully")
}

// GetNextMatchByTeamID godoc
// @Summary      Get next match for a team
// @Description  Returns the next scheduled match for a specific team
// @Tags         matches
// @ID           getNextMatchByTeamID
// @Param        id      path      int  true  "Team ID"
// @Success      200     {object}  dto.MatchResponse "Next match found"
// @Failure      400     {object}  helper.AppError "Invalid input"
// @Failure      404     {object}  helper.AppError "No upcoming matches found"
// @Failure      500     {object}  helper.AppError "Internal server error"
// @Router       /teams/{id}/next-match [get]
func (h *MatchHandler) GetNextMatchByTeamID(c *gin.Context) {
	teamIDStr := c.Param("id")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid team ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	match, err := h.MatchService.GetNextMatchByTeamID(ctx, teamID)
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}
	if match == nil {
		helper.RespondWithError(c, helper.NotFound("No upcoming matches found for this team"))
		return
	}

	matchResponse := mapper.ToMatchResponse(match)
	helper.HandleSuccess(c, http.StatusOK, matchResponse, "Next match found successfully")
}

// UpdateMatch godoc
// @Summary      Update an existing match
// @Description  Updates the details of an existing match by ID
// @Tags         matches
// @ID           updateMatch
// @Accept       json
// @Produce      json
// @Param        id     path      int                 true  "Match ID"
// @Param        match  body      dto.UpdateMatchRequest true  "Updated match data"
// @Success      200    {object}  dto.MatchResponse "Match updated"
// @Failure      400    {object}  helper.AppError "Invalid input"
// @Failure      404    {object}  helper.AppError "Match not found"
// @Failure      500    {object}  helper.AppError "Internal server error"
// @Router       /admin/matches/{id} [put]
// @Security     ApiKeyAuth
func (h *MatchHandler) UpdateMatch(c *gin.Context) {
	matchID := c.Param("id")
	id, err := strconv.ParseUint(matchID, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid match ID"))
		return
	}

	var updateRequest dto.UpdateMatchRequest
	if err = c.ShouldBindJSON(&updateRequest); err != nil {
		helper.RespondWithError(c, helper.ProcessValidationError(err, "body", "Invalid match data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	updatedMatch, err := h.MatchService.UpdateMatch(ctx, id, &updateRequest)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("match"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	matchResponse := mapper.ToMatchResponse(updatedMatch)
	helper.HandleSuccess(c, http.StatusOK, matchResponse, "Match updated successfully")
}

// DeleteMatch godoc
// @Summary      Delete a match
// @Description  Deletes a match by its ID
// @Tags         matches
// @ID           deleteMatch
// @Param        id   path      int  true  "Match ID"
// @Success      204 "No Content"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /admin/matches/{id} [delete]
// @Security     ApiKeyAuth
func (h *MatchHandler) DeleteMatch(c *gin.Context) {
	matchID := c.Param("id")
	id, err := strconv.ParseUint(matchID, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid match ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.MatchService.DeleteMatch(ctx, id)
	if err != nil && !errors.Is(err, constants.ErrRecordNotFound) {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	c.Status(http.StatusNoContent)
}
