package handler

import (
	"net/http"
	"strconv"

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
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid input", err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	role, err := h.RoleService.GetRoleByID(ctx, uint8(id))
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusNotFound, "Role not found", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, role, "Role found successfully")
}
