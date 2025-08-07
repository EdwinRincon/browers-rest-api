package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/mapper"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/service"
	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/gin-gonic/gin"
)

const (
	msgInvalidRoleID = "Invalid role ID"
)

type RoleHandler struct {
	RoleService service.RoleService
}

func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{RoleService: roleService}
}

// GetRoleByID godoc
// @Summary      Get a role by ID
// @Description  Returns the details of a role by its ID
// @Tags         roles
// @ID           getRoleByID
// @Param        id   path      int  true  "Role ID"
// @Success      200  {object}  dto.RoleResponse
// @Failure      400  {object}  helper.AppError
// @Failure      404  {object}  helper.AppError
// @Router       /roles/{id} [get]
// @Security     ApiKeyAuth
func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	roleID := c.Param("id")
	id, err := strconv.ParseUint(roleID, 10, 8)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", msgInvalidRoleID))
		return
	}

	role, err := h.RoleService.GetRoleByID(c.Request.Context(), uint8(id))
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("role"))
			return
		}
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	response := mapper.ToRoleResponse(role)
	helper.HandleSuccess(c, http.StatusOK, response, "Role retrieved successfully")
}

// CreateRole godoc
// @Summary      Create a new role
// @Description  Creates a new role with the provided data
// @Tags         roles
// @ID           createRole
// @Accept       json
// @Produce      json
// @Param        role  body      dto.CreateRoleRequest  true  "Role data"
// @Success      201   {object}  dto.RoleResponse
// @Failure      400   {object}  helper.AppError
// @Failure      409   {object}  helper.AppError
// @Router       /roles [post]
// @Security     ApiKeyAuth
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var roleDTO dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&roleDTO); err != nil {
		slog.Error("failed to bind role create request", "error", err)
		helper.RespondWithError(c, helper.BadRequest("body", "Invalid role data"))
		return
	}

	ctx := c.Request.Context()
	role := mapper.ToRole(&roleDTO)

	createdRole, err := h.RoleService.CreateRole(ctx, role)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrRecordAlreadyExists):
			helper.RespondWithError(c, helper.Conflict("role", "A role with this name already exists"))
		default:
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	response := mapper.ToRoleResponse(createdRole)
	helper.HandleSuccess(c, http.StatusCreated, response, "Role created successfully")
}

// UpdateRole godoc
// @Summary      Update an existing role
// @Description  Updates the details of an existing role by ID
// @Tags         roles
// @ID           updateRole
// @Accept       json
// @Produce      json
// @Param        id    path      int  true  "Role ID"
// @Param        role  body      dto.UpdateRoleRequest  true  "Updated role data"
// @Success      200   {object}  dto.RoleResponse
// @Failure      400   {object}  helper.AppError
// @Failure      404   {object}  helper.AppError
// @Failure      409   {object}  helper.AppError
// @Router       /roles/{id} [put]
// @Security     ApiKeyAuth
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	roleID := c.Param("id")
	id, err := strconv.ParseUint(roleID, 10, 8)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", msgInvalidRoleID))
		return
	}

	var updateRoleDTO dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&updateRoleDTO); err != nil {
		slog.Error("failed to bind role update request", "error", err)
		helper.RespondWithError(c, helper.BadRequest("body", "Invalid role data"))
		return
	}

	ctx := c.Request.Context()

	// Convert DTO to domain model
	updateRole := mapper.ToRoleFromUpdate(&updateRoleDTO)
	err = h.RoleService.UpdateRole(ctx, uint8(id), updateRole)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrRecordNotFound):
			helper.RespondWithError(c, helper.NotFound("role"))
		case errors.Is(err, constants.ErrRecordAlreadyExists):
			helper.RespondWithError(c, helper.Conflict("role", "A role with these details already exists"))
		default:
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	updatedRole, err := h.RoleService.GetActiveRoleByName(ctx, updateRole.Name)
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	response := mapper.ToRoleResponse(updatedRole)
	helper.HandleSuccess(c, http.StatusOK, response, "Role updated successfully")
}

// DeleteRole godoc
// @Summary      Delete a role
// @Description  Deletes a role by its ID
// @Tags         roles
// @ID           deleteRole
// @Param        id   path      int  true  "Role ID"
// @Success      204 "No Content"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Router       /roles/{id} [delete]
// @Security     ApiKeyAuth
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	roleID := c.Param("id")
	id, err := strconv.ParseUint(roleID, 10, 8)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", msgInvalidRoleID))
		return
	}

	ctx := c.Request.Context()
	_ = h.RoleService.DeleteRole(ctx, uint8(id))

	c.Status(http.StatusNoContent)
}

// GetPaginatedRoles godoc
// @Summary      List paginated roles
// @Description  Retrieves a paginated list of roles with optional sorting
// @Tags         roles
// @ID           listPaginatedRoles
// @Produce      json
// @Param        sort      query     string  false  "Field to sort by (e.g. name, created_at)"
// @Param        order     query     string  false  "Sort direction: asc or desc"  Enums(asc, desc)
// @Param        page      query     int     false  "Page number (starts at 0)"
// @Param        pageSize  query     int     false  "Number of items per page"
// @Success      200       {object}  helper.PaginatedResponse{items=[]dto.RoleResponse, totalCount=int}
// @Failure      400       {object}  helper.AppError "Invalid input"
// @Failure      500       {object}  helper.AppError
// @Router       /roles [get]
// @Security     ApiKeyAuth
func (h *RoleHandler) GetPaginatedRoles(c *gin.Context) {
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 0 {
		page = 1
	}
	if pageSize < 0 || pageSize > 100 {
		pageSize = 10
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	// Validate sort field
	if err := helper.ValidateSort(model.Role{}, sort); err != nil {
		helper.RespondWithError(c, helper.BadRequest("sort", err.Error()))
		return
	}

	ctx := c.Request.Context()
	roles, total, err := h.RoleService.GetPaginatedRoles(ctx, sort, order, page, pageSize)
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	response := helper.PaginatedResponse{
		Items:      mapper.ToRoleResponseList(roles),
		TotalCount: total,
	}

	helper.HandleSuccess(c, http.StatusOK, response, "Roles retrieved successfully")
}
