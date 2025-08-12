package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/service"
	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TeamStatsHandler struct {
	TeamStatsService service.TeamStatsService
}

func NewTeamStatsHandler(teamStatsService service.TeamStatsService) *TeamStatsHandler {
	return &TeamStatsHandler{TeamStatsService: teamStatsService}
}

// CreateTeamStats godoc
// @Summary      Create new team stats
// @Description  Creates new team stats with the provided data
// @Tags         team-stats
// @ID           createTeamStats
// @Accept       json
// @Produce      json
// @Param        teamStats  body      model.TeamStat  true  "Team stats data"
// @Success      201        {object}  model.TeamStat "Created"
// @Failure      400        {object}  helper.AppError "Invalid input"
// @Failure      409        {object}  helper.AppError "Conflict (e.g., stats for this team/season already exist)"
// @Failure      500        {object}  helper.AppError "Internal server error"
// @Router       /team-stats [post]
// @Security     ApiKeyAuth
func (h *TeamStatsHandler) CreateTeamStats(c *gin.Context) {
	var teamStats model.TeamStat
	if err := c.ShouldBindJSON(&teamStats); err != nil {
		helper.RespondWithError(c, helper.BadRequest("body", "Invalid team stats data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	err := h.TeamStatsService.CreateTeamStats(ctx, &teamStats)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			helper.RespondWithError(c, helper.Conflict("team-stats", "Stats for this team/season already exist"))
			return
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
			return
		}
	}

	helper.HandleSuccess(c, http.StatusCreated, teamStats, "Team stats created successfully")
}

// GetTeamStatsByID godoc
// @Summary      Get team stats by ID
// @Description  Returns the details of team stats by its ID
// @Tags         team-stats
// @ID           getTeamStatsByID
// @Param        id   path      int  true  "Team Stats ID"
// @Success      200  {object}  model.TeamStat "Success"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Team stats not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /team-stats/{id} [get]
// @Security     ApiKeyAuth
func (h *TeamStatsHandler) GetTeamStatsByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid team stats ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	teamStats, err := h.TeamStatsService.GetTeamStatsByID(ctx, id)
	if err != nil {
		helper.RespondWithError(c, helper.NotFound("team stats"))
		return
	}

	helper.HandleSuccess(c, http.StatusOK, teamStats, "Team stats retrieved successfully")
}

// ListTeamStats godoc
// @Summary      List team stats
// @Description  Retrieves a paginated list of team stats
// @Tags         team-stats
// @ID           listTeamStats
// @Param        page  query     int  false  "Page number"
// @Success      200   {array}   model.TeamStat "Team stats listed successfully"
// @Failure      400   {object}  helper.AppError "Invalid input"
// @Failure      500   {object}  helper.AppError "Internal server error"
// @Router       /team-stats [get]
// @Security     ApiKeyAuth
func (h *TeamStatsHandler) ListTeamStats(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("page", "Invalid page number"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	teamStats, err := h.TeamStatsService.ListTeamStats(ctx, page)
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	helper.HandleSuccess(c, http.StatusOK, teamStats, "Team stats listed successfully")
}

// UpdateTeamStats godoc
// @Summary      Update existing team stats
// @Description  Updates the details of existing team stats by ID
// @Tags         team-stats
// @ID           updateTeamStats
// @Accept       json
// @Produce      json
// @Param        id         path      int             true  "Team Stats ID"
// @Param        teamStats  body      model.TeamStat  true  "Updated team stats data"
// @Success      200        {object}  model.TeamStat "Updated"
// @Failure      400        {object}  helper.AppError "Invalid input"
// @Failure      404        {object}  helper.AppError "Team stats not found"
// @Failure      500        {object}  helper.AppError "Internal server error"
// @Router       /team-stats/{id} [put]
// @Security     ApiKeyAuth
func (h *TeamStatsHandler) UpdateTeamStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid team stats ID"))
		return
	}

	var teamStats model.TeamStat
	if err = c.ShouldBindJSON(&teamStats); err != nil {
		helper.RespondWithError(c, helper.BadRequest("body", "Invalid team stats data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	teamStats.ID = id
	err = h.TeamStatsService.UpdateTeamStats(ctx, &teamStats)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("team stats"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusOK, teamStats, "Team stats updated successfully")
}

// DeleteTeamStats godoc
// @Summary      Delete team stats
// @Description  Deletes team stats by its ID
// @Tags         team-stats
// @ID           deleteTeamStats
// @Param        id   path      int  true  "Team Stats ID"
// @Success      204 "No Content"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Team stats not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /team-stats/{id} [delete]
// @Security     ApiKeyAuth
func (h *TeamStatsHandler) DeleteTeamStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid team stats ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	err = h.TeamStatsService.DeleteTeamStats(ctx, id)
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	c.Status(http.StatusNoContent)
}
