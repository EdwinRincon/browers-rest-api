package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	httpMapper "github.com/EdwinRincon/browersfc-api/adapter/http"
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/EdwinRincon/browersfc-api/internal/domain/service"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	MatchDomainService *service.MatchDomainService
	MatchMapper        *httpMapper.MatchHTTPMapper
}

func NewMatchHandler(matchDomainService *service.MatchDomainService) *MatchHandler {
	return &MatchHandler{
		MatchDomainService: matchDomainService,
		MatchMapper:        httpMapper.NewMatchHTTPMapper(),
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
// @Security     BearerAuth
func (h *MatchHandler) CreateMatch(c *gin.Context) {
	var createRequest dto.CreateMatchRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		helper.WriteErrorResponse(c, helper.BuildValidationErrorFromBinding(err, "body", "Invalid match data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Convert DTO to domain entity
	domainMatch := h.MatchMapper.DTOToDomain(&createRequest)

	createdMatch, err := h.MatchDomainService.CreateMatch(ctx, domainMatch)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	matchResponse := h.MatchMapper.DomainToShortDTO(createdMatch)
	helper.WriteSuccessResponse(c, http.StatusCreated, matchResponse, "Match created successfully")
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
// @Security     BearerAuth
func (h *MatchHandler) GetMatchByID(c *gin.Context) {
	matchID := c.Param("id")
	id, err := strconv.ParseUint(matchID, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid match ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	match, err := h.MatchDomainService.GetMatchByID(ctx, id)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("match"))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	matchResponse := h.MatchMapper.DomainToDTO(match)
	helper.WriteSuccessResponse(c, http.StatusOK, matchResponse, "Match found successfully")
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
// @Security     BearerAuth
func (h *MatchHandler) GetDetailedMatchByID(c *gin.Context) {
	matchID := c.Param("id")
	id, err := strconv.ParseUint(matchID, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid match ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	match, err := h.MatchDomainService.GetDetailedMatchByID(ctx, id)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("match"))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	detailedResponse := h.MatchMapper.DomainToDetailDTO(match)
	helper.WriteSuccessResponse(c, http.StatusOK, detailedResponse, "Match details found successfully")
}

// GetPaginatedMatches godoc
// @Summary      Get paginated matches
// @Description  Retrieves a paginated list of matches with sorting and ordering
// @Summary Get paginated matches
// @Description Retrieve matches with pagination, optional sorting, and ordering.
// @Tags matches
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(0)
// @Param pageSize query int false "Page size" default(10)
// @Param sort query string false "Sort field"
// @Param order query string false "Sort order" Enums(asc, desc) default(desc)
// @Success 200 {object} helper.AppSuccess{data=helper.PaginatedResponse{items=[]dto.MatchResponse, totalCount=int}}
// @Failure 400 {object} helper.AppError "Invalid input"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Router /matches [get]
func (h *MatchHandler) GetPaginatedMatches(c *gin.Context) {
	sort := c.DefaultQuery("sort", "kickoff")
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
		helper.WriteErrorResponse(c, helper.NewBadRequestError("sort", err.Error()))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	matches, total, err := h.MatchDomainService.GetPaginatedMatches(ctx, sort, order, page, pageSize)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	response := helper.PaginatedResponse{
		Items:      h.MatchMapper.DomainListToDTO(matches),
		TotalCount: total,
	}

	helper.WriteSuccessResponse(c, http.StatusOK, response, "Matches retrieved successfully")
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
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid season ID"))
		return
	}

	sort := c.DefaultQuery("sort", "kickoff")
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

	matches, total, err := h.MatchDomainService.GetMatchesBySeasonID(ctx, seasonID, sort, order, page, pageSize)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	response := helper.PaginatedResponse{
		Items:      h.MatchMapper.DomainListToDTO(matches),
		TotalCount: total,
	}

	helper.WriteSuccessResponse(c, http.StatusOK, response, "Season matches retrieved successfully")
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
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid team ID"))
		return
	}

	sort := c.DefaultQuery("sort", "kickoff")
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

	matches, total, err := h.MatchDomainService.GetMatchesByTeamID(ctx, teamID, sort, order, page, pageSize)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	response := helper.PaginatedResponse{
		Items:      h.MatchMapper.DomainListToDTO(matches),
		TotalCount: total,
	}

	helper.WriteSuccessResponse(c, http.StatusOK, response, "Team matches retrieved successfully")
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
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid team ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	match, err := h.MatchDomainService.GetNextMatchByTeamID(ctx, teamID)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}
	if match == nil {
		helper.WriteErrorResponse(c, helper.NewNotFoundError("No upcoming matches found for this team"))
		return
	}

	matchResponse := h.MatchMapper.DomainToDTO(match)
	helper.WriteSuccessResponse(c, http.StatusOK, matchResponse, "Next match found successfully")
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
// @Security     BearerAuth
func (h *MatchHandler) UpdateMatch(c *gin.Context) {
	matchID := c.Param("id")
	id, err := strconv.ParseUint(matchID, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid match ID"))
		return
	}

	var updateRequest dto.UpdateMatchRequest
	if err = c.ShouldBindJSON(&updateRequest); err != nil {
		helper.WriteErrorResponse(c, helper.BuildValidationErrorFromBinding(err, "body", "Invalid match data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Convert DTO to domain entity
	domainMatch := h.MatchMapper.UpdateDTOToDomain(&updateRequest)

	updatedMatch, err := h.MatchDomainService.UpdateMatch(ctx, id, domainMatch)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("match"))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	matchResponse := h.MatchMapper.DomainToDTO(updatedMatch)
	helper.WriteSuccessResponse(c, http.StatusOK, matchResponse, "Match updated successfully")
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
// @Security     BearerAuth
func (h *MatchHandler) DeleteMatch(c *gin.Context) {
	matchID := c.Param("id")
	id, err := strconv.ParseUint(matchID, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid match ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.MatchDomainService.DeleteMatch(ctx, id)
	if err != nil && !errors.Is(err, constants.ErrRecordNotFound) {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	c.Status(http.StatusNoContent)
}
