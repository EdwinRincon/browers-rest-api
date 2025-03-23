package handler

import (
	"net/http"
	"strconv"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/service"
	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	RoleService service.RoleService
}

func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{
		RoleService: roleService,
	}
}

func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	roleID := c.Param("id")
	id, err := strconv.ParseUint(roleID, 10, 8)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), false)
		return
	}

	ctx := c.Request.Context()
	role, err := h.RoleService.GetRoleByID(ctx, uint8(id))
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, role, "Role retrieved successfully")
}

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

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	roleID := c.Param("id")
	id, err := strconv.ParseUint(roleID, 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		helper.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	role.ID = uint8(id)
	err = h.RoleService.UpdateRole(ctx, &role)
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, role, "Role updated successfully")
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	roleID := c.Param("id")
	id, err := strconv.ParseUint(roleID, 10, 8)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err = h.RoleService.DeleteRole(ctx, uint8(id))
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, nil, "Role deleted successfully")
}

func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	ctx := c.Request.Context()
	roles, err := h.RoleService.GetAllRoles(ctx)
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, roles, "Role retrieved successfully")
}
