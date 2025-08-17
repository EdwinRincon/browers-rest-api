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
	"gorm.io/gorm"
)

type TeamHandler struct {
	TeamService service.TeamService
}

func NewTeamHandler(teamService service.TeamService) *TeamHandler {
	return &TeamHandler{
		TeamService: teamService,
	}
}

// CreateTeam godoc
// @Summary      Create a new team
// @Description  Creates a new team with the provided data
// @Tags         teams
// @ID           createTeam
// @Accept       json
// @Produce      json
// @Param        team  body      dto.CreateTeamRequest  true  "Team data"
// @Success      201   {object}  dto.TeamShort "Created"
// @Failure      400   {object}  helper.AppError "Invalid input"
// @Failure      409   {object}  helper.AppError "Conflict (e.g., team name exists)"
// @Failure      500   {object}  helper.AppError "Internal server error"
// @Router       /teams [post]
// @Security     ApiKeyAuth
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var createRequest dto.CreateTeamRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		helper.RespondWithError(c, helper.BadRequest("body", "Invalid team data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Map DTO to model
	team := mapper.ToTeam(&createRequest)

	createdTeam, err := h.TeamService.CreateTeam(ctx, team)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrRecordAlreadyExists):
			helper.RespondWithError(c, helper.Conflict("team", "A team with this name already exists"))
			return
		default:
			helper.RespondWithError(c, helper.InternalError(err))
			return
		}
	}

	helper.HandleSuccess(c, http.StatusCreated, createdTeam, "Team created successfully")
}

// GetTeamByID godoc
// @Summary      Get a team by ID
// @Description  Returns the details of a team by its ID
// @Tags         teams
// @ID           getTeamByID
// @Param        id   path      int  true  "Team ID"
// @Success      200  {object}  model.Team "Success"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Team not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /teams/{id} [get]
// @Security     ApiKeyAuth
func (h *TeamHandler) GetTeamByID(c *gin.Context) {
	teamID := c.Param("id")
	id, err := strconv.ParseUint(teamID, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid team ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	team, err := h.TeamService.GetTeamByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("team"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusOK, team, "Team found successfully")
}

// GetPaginatedTeams godoc
// @Summary      List teams
// @Description  Returns a paginated list of teams
// @Tags         teams
// @ID           GetPaginatedTeams
// @Param        page      query     int     false  "Page number (default: 0)"
// @Param        pageSize  query     int     false  "Page size (default: 10)"
// @Param        sort      query     string  false  "Sort field (default: id)"
// @Param        order     query     string  false  "Sort order (asc or desc, default: asc)"
// @Success      200       {object}  helper.PaginatedResponse{data=[]dto.TeamResponse} "Success"
// @Failure      400       {object}  helper.AppError "Invalid input"
// @Failure      500       {object}  helper.AppError "Internal server error"
// @Router       /teams [get]
// @Security     ApiKeyAuth
func (h *TeamHandler) GetPaginatedTeams(c *gin.Context) {
	sort := c.DefaultQuery("sort", "id")
	order := c.DefaultQuery("order", "asc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 0 {
		page = 0
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Validate sort field
	if err := helper.ValidateSort(model.Team{}, sort); err != nil {
		helper.RespondWithError(c, helper.BadRequest("sort", err.Error()))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	teams, total, err := h.TeamService.GetPaginatedTeams(ctx, sort, order, page, pageSize)
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	response := helper.PaginatedResponse{
		Items:      mapper.ToTeamResponseList(teams),
		TotalCount: total,
	}

	helper.HandleSuccess(c, http.StatusOK, response, "Teams retrieved successfully")
}

// UpdateTeam godoc
// @Summary      Update a team
// @Description  Updates a team with the provided data
// @Tags         teams
// @ID           updateTeam
// @Accept       json
// @Produce      json
// @Param        id    path      int                  true  "Team ID"
// @Param        team  body      dto.UpdateTeamRequest  true  "Updated team data"
// @Success      200   {object}  dto.TeamResponse "Success"
// @Failure      400   {object}  helper.AppError "Invalid input"
// @Failure      404   {object}  helper.AppError "Team not found"
// @Failure      409   {object}  helper.AppError "Conflict (e.g., team name exists)"
// @Failure      500   {object}  helper.AppError "Internal server error"
// @Router       /teams/{id} [put]
// @Security     ApiKeyAuth
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	teamID := c.Param("id")
	id, err := strconv.ParseUint(teamID, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid team ID"))
		return
	}

	var updateTeamRequest dto.UpdateTeamRequest
	if err = c.ShouldBindJSON(&updateTeamRequest); err != nil {
		helper.RespondWithError(c, helper.BadRequest("body", "Invalid team data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	updatedTeam, err := h.TeamService.UpdateTeam(ctx, &updateTeamRequest, id)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("team"))
		} else if errors.Is(err, constants.ErrRecordAlreadyExists) {
			helper.RespondWithError(c, helper.Conflict("team", "A team with this name already exists"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	teamResponse := mapper.ToTeamResponse(updatedTeam)
	helper.HandleSuccess(c, http.StatusOK, teamResponse, "Team updated successfully")
}

// DeleteTeam godoc
// @Summary      Delete a team
// @Description  Deletes a team by its ID
// @Tags         teams
// @ID           deleteTeam
// @Param        id   path      int  true  "Team ID"
// @Success      204 "No Content"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Team not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /teams/{id} [delete]
// @Security     ApiKeyAuth
func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	teamID := c.Param("id")
	id, err := strconv.ParseUint(teamID, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid team ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	err = h.TeamService.DeleteTeam(ctx, id)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("team"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	c.Status(http.StatusNoContent)
}
