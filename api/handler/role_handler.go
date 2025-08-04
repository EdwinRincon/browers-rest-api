package handler

import (
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

type RoleHandler struct {
	RoleService service.RoleService
}

func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{
		RoleService: roleService,
	}
}

// GetRoleByID godoc
// @Summary      Get a role by ID
// @Description  Returns the details of a role by its ID
// @Tags         roles
// @ID           getRoleByID
// @Param        id   path      int  true  "Role ID"
// @Success      200  {object}  model.Role  "Success"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Role not found"
// @Router       /roles/{id} [get]
// @Security     ApiKeyAuth
func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	roleID := c.Param("id")
	id, err := strconv.ParseUint(roleID, 10, 8)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid role ID"))
		return
	}

	ctx := c.Request.Context()
	role, err := h.RoleService.GetRoleByID(ctx, uint8(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("role"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusOK, role, "Role retrieved successfully")
}

// CreateRole godoc
// @Summary      Create a new role
// @Description  Creates a new role with the provided data
// @Tags         roles
// @ID           createRole
// @Accept       json
// @Produce      json
// @Param        role  body      model.Role  true  "Role data"
// @Success      201   {object}  model.Role  "Created"
// @Failure      400   {object}  helper.AppError "Invalid input"
// @Failure      409   {object}  helper.AppError "Role already exists"
// @Router       /roles [post]
// @Security     ApiKeyAuth
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		helper.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.RoleService.CreateRole(ctx, &role)
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusCreated, role, "Role created successfully")
}

// UpdateRole godoc
// @Summary      Update an existing role
// @Description  Updates the details of an existing role by ID
// @Tags         roles
// @ID           updateRole
// @Accept       json
// @Produce      json
// @Param        id    path      int        true  "Role ID"
// @Param        role  body      model.Role true  "Updated role data"
// @Success      200   {object}  model.Role  "Updated"
// @Failure      400   {object}  helper.AppError "Invalid input"
// @Failure      404   {object}  helper.AppError "Role not found"
// @Router       /roles/{id} [put]
// @Security     ApiKeyAuth
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	roleID := c.Param("id")
	id, err := strconv.ParseUint(roleID, 10, 8)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid role ID"))
		return
	}

	ctx := c.Request.Context()
	existingRole, err := h.RoleService.GetRoleByID(ctx, uint8(id))
	if err != nil {
		helper.RespondWithError(c, helper.NotFound("role"))
		return
	}

	var updateRole model.Role
	if err := c.ShouldBindJSON(&updateRole); err != nil {
		helper.HandleValidationError(c, err)
		return
	}

	updateRole.ID = existingRole.ID
	updateRole.CreatedAt = existingRole.CreatedAt
	updateRole.UpdatedAt = time.Now()

	if id > 255 {
		helper.RespondWithError(c, helper.BadRequest("id", "Role ID too large"))
		return
	}

	err = h.RoleService.UpdateRole(ctx, &updateRole)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("role"))
		} else if errors.Is(err, gorm.ErrDuplicatedKey) {
			helper.RespondWithError(c, helper.Conflict("role", "A role with these details already exists"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusOK, updateRole, "Role updated successfully")
}

// DeleteRole godoc
// @Summary      Delete a role
// @Description  Deletes a role by its ID
// @Tags         roles
// @ID           deleteRole
// @Param        id   path      int  true  "Role ID"
// @Success      204  {object}  nil  "No Content"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Role not found"
// @Router       /roles/{id} [delete]
// @Security     ApiKeyAuth
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	roleID := c.Param("id")
	id, err := strconv.ParseUint(roleID, 10, 8)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid role ID"))
		return
	}

	ctx := c.Request.Context()
	err = h.RoleService.DeleteRole(ctx, uint8(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("role"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAllRoles godoc
// @Summary      List roles
// @Description  Retrieves a list of all roles
// @Tags         roles
// @ID           listRoles
// @Success      200  {array}   model.Role  "Success"
// @Router       /roles [get]
// @Security     ApiKeyAuth
func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	ctx := c.Request.Context()
	roles, err := h.RoleService.GetAllRoles(ctx)
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	helper.HandleSuccess(c, http.StatusOK, roles, "Role retrieved successfully")
}
