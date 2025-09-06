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
	"github.com/EdwinRincon/browersfc-api/helper"
	domainservice "github.com/EdwinRincon/browersfc-api/internal/domain/service"
	"github.com/gin-gonic/gin"
)

const (
	msgInvalidRoleID = "Invalid role ID"
)

type RoleHandler struct {
	RoleDomainService *domainservice.RoleDomainService
}

func NewRoleHandler(roleDomainService *domainservice.RoleDomainService) *RoleHandler {
	return &RoleHandler{RoleDomainService: roleDomainService}
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
// @Router       /admin/roles/{id} [get]
// @Security     ApiKeyAuth
func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	roleID := c.Param("id")
	id, err := strconv.ParseUint(roleID, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", msgInvalidRoleID))
		return
	}

	domainRole, err := h.RoleDomainService.GetRoleByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domainservice.ErrRoleNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("role"))
			return
		}
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	// Convert domain role directly to response DTO
	response := mapper.DomainToRoleResponse(domainRole)
	helper.WriteSuccessResponse(c, http.StatusOK, response, "Role retrieved successfully")
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
// @Router       /admin/roles [post]
// @Security     ApiKeyAuth
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var roleDTO dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&roleDTO); err != nil {
		helper.WriteErrorResponse(c, helper.BuildValidationErrorFromBinding(err, "body", constants.MsgInvalidRoleData))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Convert DTO to domain role
	roleMapper := mapper.NewRoleDomainMapper()
	modelRole := mapper.ToRole(&roleDTO)
	domainRole := roleMapper.ToDomain(modelRole)

	createdRole, err := h.RoleDomainService.CreateRole(ctx, domainRole)
	if err != nil {
		switch {
		case errors.Is(err, domainservice.ErrRoleAlreadyExists):
			helper.WriteErrorResponse(c, helper.NewConflictError("role", "A role with this name already exists"))
		case errors.Is(err, domainservice.ErrInvalidRole):
			helper.WriteErrorResponse(c, helper.NewBadRequestError("role", "Invalid role data"))
		default:
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	// Convert domain role directly to response DTO
	response := mapper.DomainToRoleResponse(createdRole)
	helper.WriteSuccessResponse(c, http.StatusCreated, response, "Role created successfully")
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
// @Router       /admin/roles/{id} [put]
// @Security     ApiKeyAuth
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	roleID := c.Param("id")
	id, err := strconv.ParseUint(roleID, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", msgInvalidRoleID))
		return
	}

	var updateRoleDTO dto.UpdateRoleRequest
	if err = c.ShouldBindJSON(&updateRoleDTO); err != nil {
		helper.WriteErrorResponse(c, helper.BuildValidationErrorFromBinding(err, "body", constants.MsgInvalidRoleData))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Convert DTO to domain model
	roleMapper := mapper.NewRoleDomainMapper()
	updateModel := mapper.ToRoleFromUpdate(&updateRoleDTO)
	updateDomain := roleMapper.ToDomain(updateModel)
	updateDomain.ID = id // Set the ID from URL parameter

	updatedRole, err := h.RoleDomainService.UpdateRole(ctx, updateDomain)
	if err != nil {
		switch {
		case errors.Is(err, domainservice.ErrRoleNotFound):
			helper.WriteErrorResponse(c, helper.NewNotFoundError("role"))
		case errors.Is(err, domainservice.ErrRoleAlreadyExists):
			helper.WriteErrorResponse(c, helper.NewConflictError("role", "A role with these details already exists"))
		case errors.Is(err, domainservice.ErrInvalidRole):
			helper.WriteErrorResponse(c, helper.NewBadRequestError("role", "Invalid role data"))
		default:
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	// Convert domain role directly to response DTO
	response := mapper.DomainToRoleResponse(updatedRole)
	helper.WriteSuccessResponse(c, http.StatusOK, response, "Role updated successfully")
}

// DeleteRole godoc
// @Summary      Delete a role
// @Description  Deletes a role by its ID
// @Tags         roles
// @ID           deleteRole
// @Param        id   path      int  true  "Role ID"
// @Success      204 "No Content"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Router       /admin/roles/{id} [delete]
// @Security     ApiKeyAuth
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	roleID := c.Param("id")
	id, err := strconv.ParseUint(roleID, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", msgInvalidRoleID))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.RoleDomainService.DeleteRole(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, domainservice.ErrRoleNotFound):
			helper.WriteErrorResponse(c, helper.NewNotFoundError("role"))
		case errors.Is(err, domainservice.ErrCannotDeleteSystemRole):
			helper.WriteErrorResponse(c, helper.NewBadRequestError("role", "Cannot delete system role"))
		case errors.Is(err, domainservice.ErrInvalidRole):
			helper.WriteErrorResponse(c, helper.NewBadRequestError("id", msgInvalidRoleID))
		default:
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

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
// @Router       /admin/roles [get]
// @Security     ApiKeyAuth
func (h *RoleHandler) GetPaginatedRoles(c *gin.Context) {
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 0 {
		page = 0
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	// Validate sort field
	if err := helper.ValidateSort(model.Role{}, sort); err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("sort", err.Error()))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Convert page to 1-based indexing for domain service
	domainRoles, total, err := h.RoleDomainService.GetPaginatedRoles(ctx, sort, order, page, pageSize)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	// Convert domain roles directly to response DTOs
	response := helper.PaginatedResponse{
		Items:      mapper.DomainToRoleResponseList(domainRoles),
		TotalCount: total,
	}

	helper.WriteSuccessResponse(c, http.StatusOK, response, "Roles retrieved successfully")
}
