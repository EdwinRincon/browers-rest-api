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
// @Param        team  body      model.Team  true  "Team data"
// @Success      201   {object}  model.Team "Created"
// @Failure      400   {object}  helper.AppError "Invalid input"
// @Failure      409   {object}  helper.AppError "Conflict (e.g., team name exists)"
// @Failure      500   {object}  helper.AppError "Internal server error"
// @Router       /teams [post]
// @Security     ApiKeyAuth
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var team model.Team
	if err := c.ShouldBindJSON(&team); err != nil {
		helper.RespondWithError(c, helper.BadRequest("body", "Invalid team data"))
		return
	}

	// Reset ID to ensure we're creating a new record
	team.ID = 0

	ctx := c.Request.Context()
	err := h.TeamService.CreateTeam(ctx, &team)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			helper.RespondWithError(c, helper.Conflict("team", "A team with this name already exists"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusCreated, team, "Team created successfully")
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

	ctx := c.Request.Context()
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

// ListTeams godoc
// @Summary      List teams
// @Description  Retrieves a paginated list of teams
// @Tags         teams
// @ID           listTeams
// @Param        page  query     int  false  "Page number"
// @Success      200   {array}   model.Team "Teams listed successfully"
// @Failure      400   {object}  helper.AppError "Invalid input"
// @Failure      500   {object}  helper.AppError "Internal server error"
// @Router       /teams [get]
// @Security     ApiKeyAuth
func (h *TeamHandler) ListTeams(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("page", "Invalid page number"))
		return
	}

	ctx := c.Request.Context()
	teams, err := h.TeamService.ListTeams(ctx, page)
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	helper.HandleSuccess(c, http.StatusOK, teams, "Teams listed successfully")
}

// UpdateTeam godoc
// @Summary      Update an existing team
// @Description  Updates the details of an existing team by ID
// @Tags         teams
// @ID           updateTeam
// @Accept       json
// @Produce      json
// @Param        id    path      int        true  "Team ID"
// @Param        team  body      model.Team true  "Updated team data"
// @Success      200   {object}  model.Team "Updated"
// @Failure      400   {object}  helper.AppError "Invalid input"
// @Failure      404   {object}  helper.AppError "Team not found"
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

	var team model.Team
	if err := c.ShouldBindJSON(&team); err != nil {
		helper.RespondWithError(c, helper.BadRequest("body", "Invalid team data"))
		return
	}

	ctx := c.Request.Context()
	team.ID = id
	err = h.TeamService.UpdateTeam(ctx, &team)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("team"))
		} else if errors.Is(err, gorm.ErrDuplicatedKey) {
			helper.RespondWithError(c, helper.Conflict("team", "A team with this name already exists"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusOK, team, "Team updated successfully")
}

// DeleteTeam godoc
// @Summary      Delete a team
// @Description  Deletes a team by its ID
// @Tags         teams
// @ID           deleteTeam
// @Param        id   path      int  true  "Team ID"
// @Success      204  {object}  nil "No Content"
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

	ctx := c.Request.Context()
	err = h.TeamService.DeleteTeam(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("team"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	c.Status(http.StatusNoContent)
}
