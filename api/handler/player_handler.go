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

type PlayerHandler struct {
	PlayerService service.PlayerService
}

func NewPlayerHandler(playerService service.PlayerService) *PlayerHandler {
	return &PlayerHandler{
		PlayerService: playerService,
	}
}

// CreatePlayer godoc
// @Summary      Create a new player
// @Description  Creates a new player with the provided data
// @Tags         players
// @ID           createPlayer
// @Accept       json
// @Produce      json
// @Param        player  body      model.Player  true  "Player data"
// @Success      201     {object}  model.Player  "Created"
// @Failure      400     {object}  helper.AppError "Invalid input"
// @Failure      409     {object}  helper.AppError "Player already exists"
// @Router       /players [post]
// @Security     ApiKeyAuth
func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var player model.Player
	if err := c.ShouldBindJSON(&player); err != nil {
		helper.RespondWithError(c, helper.BadRequest("body", "Invalid player data"))
		return
	}

	player.ID = 0

	ctx := c.Request.Context()
	err := h.PlayerService.CreatePlayer(ctx, &player)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			helper.RespondWithError(c, helper.Conflict("player", "A player with these details already exists"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusCreated, player, "Player created successfully")
}

// GetPlayerByID godoc
// @Summary      Get a player by ID
// @Description  Returns the details of a player by its ID
// @Tags         players
// @ID           getPlayerByID
// @Param        id   path      int  true  "Player ID"
// @Success      200  {object}  model.Player  "Success"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Player not found"
// @Router       /players/{id} [get]
// @Security     ApiKeyAuth
func (h *PlayerHandler) GetPlayerByID(c *gin.Context) {
	playerID := c.Param("id")
	id, err := strconv.ParseUint(playerID, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid player ID"))
		return
	}

	ctx := c.Request.Context()
	player, err := h.PlayerService.GetPlayerByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("player"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusOK, player, "Player found successfully")
}

// UpdatePlayer godoc
// @Summary      Update an existing player
// @Description  Updates the details of an existing player by ID
// @Tags         players
// @ID           updatePlayer
// @Accept       json
// @Produce      json
// @Param        id      path      int         true  "Player ID"
// @Param        player  body      model.Player true  "Updated player data"
// @Success      200     {object}  model.Player  "Updated"
// @Failure      400     {object}  helper.AppError "Invalid input"
// @Failure      404     {object}  helper.AppError "Player not found"
// @Router       /players/{id} [put]
// @Security     ApiKeyAuth
func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid player ID"))
		return
	}

	var player model.Player
	if err := c.ShouldBindJSON(&player); err != nil {
		helper.RespondWithError(c, helper.BadRequest("body", "Invalid player data"))
		return
	}

	player.ID = id

	ctx := c.Request.Context()
	err = h.PlayerService.UpdatePlayer(ctx, &player)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("player"))
		} else if errors.Is(err, gorm.ErrDuplicatedKey) {
			helper.RespondWithError(c, helper.Conflict("player", "A player with these details already exists"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusOK, player, "Player updated successfully")
}

// DeletePlayer godoc
// @Summary      Delete a player
// @Description  Deletes a player by its ID
// @Tags         players
// @ID           deletePlayer
// @Param        id   path      int  true  "Player ID"
// @Success      204  {object}  nil  "No Content"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Player not found"
// @Router       /players/{id} [delete]
// @Security     ApiKeyAuth
func (h *PlayerHandler) DeletePlayer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid player ID"))
		return
	}

	ctx := c.Request.Context()
	err = h.PlayerService.DeletePlayer(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("player"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusOK, nil, "Player deleted successfully")
}

// GetAllPlayers godoc
// @Summary      List players
// @Description  Retrieves a paginated list of players
// @Tags         players
// @ID           listPlayers
// @Param        page  query     int  false  "Page number"
// @Success      200   {array}   model.Player  "Success"
// @Failure      400   {object}  helper.AppError "Invalid input"
// @Router       /players [get]
// @Security     ApiKeyAuth
func (h *PlayerHandler) GetAllPlayers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("page", "Invalid page number"))
		return
	}

	ctx := c.Request.Context()
	players, err := h.PlayerService.GetAllPlayers(ctx, page)
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}
	helper.HandleSuccess(c, http.StatusOK, players, "Players retrieved successfully")
}
