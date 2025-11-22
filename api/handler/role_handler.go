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
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/helper"
	domainservice "github.com/EdwinRincon/browersfc-api/internal/domain/service"
	"github.com/gin-gonic/gin"
)

const (
	msgInvalidRoleID = "Invalid role ID"
)

type RoleHandler struct {
	RoleDomainService *domainservice.RoleDomainService
	RoleMapper        *httpMapper.RoleHTTPMapper
}

func NewRoleHandler(roleDomainService *domainservice.RoleDomainService) *RoleHandler {
	return &RoleHandler{
		RoleDomainService: roleDomainService,
		RoleMapper:        httpMapper.NewRoleHTTPMapper(),
	}
}

// GetRoleByID godoc
// @Summary      Get a role by ID
// @Description  Returns the details of a role by its ID
// @Tags         roles
// @ID           getRoleByID
// @Param        id   path      int  true  "Role ID"
// @Success      200  {object}  dto.RoleResponse
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Role not found"
// @Router       /admin/roles/{id} [get]
// @Security     BearerAuth
func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	roleID := c.Param("id")
	id, err := strconv.ParseUint(roleID, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", msgInvalidRoleID))
		return
	}

	domainRole, err := h.RoleDomainService.GetRoleByID(c.Request.Context(), id)
	if err != nil {
		if err == constants.ErrRecordNotFound {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("role"))
			return
		}
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	// Convert domain role directly to response DTO
	response := h.RoleMapper.DomainToDTO(domainRole)
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
// @Failure      400   {object}  helper.AppError "Invalid input"
// @Failure      409   {object}  helper.AppError "Role already exists"
// @Router       /admin/roles [post]
// @Security     BearerAuth
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
	domainRole := h.RoleMapper.DTOToDomain(&roleDTO)

	createdRole, err := h.RoleDomainService.CreateRole(ctx, domainRole)
	if err != nil {
		if err == constants.ErrRecordAlreadyExists {
			helper.WriteErrorResponse(c, helper.NewConflictError("role", "A role with this name already exists"))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	// Convert domain role directly to response DTO
	response := h.RoleMapper.DomainToDTO(createdRole)
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
// @Failure      400   {object}  helper.AppError "Invalid input"
// @Failure      404   {object}  helper.AppError "Role not found"
// @Failure      409   {object}  helper.AppError "Role already exists"
// @Router       /admin/roles/{id} [put]
// @Security     BearerAuth
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
	updateDomain := h.RoleMapper.UpdateDTOToDomain(&updateRoleDTO)
	updateDomain.ID = id // Set the ID from URL parameter

	updatedRole, err := h.RoleDomainService.UpdateRole(ctx, updateDomain)
	if err != nil {
		switch err {
		case constants.ErrRecordNotFound:
			helper.WriteErrorResponse(c, helper.NewNotFoundError("role"))
		case constants.ErrRecordAlreadyExists:
			helper.WriteErrorResponse(c, helper.NewConflictError("role", "A role with these details already exists"))
		default:
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	// Convert domain role directly to response DTO
	response := h.RoleMapper.DomainToDTO(updatedRole)
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
// @Security     BearerAuth
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
		if err == constants.ErrRecordNotFound {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("role"))
		} else if errors.Is(err, constants.ErrCannotDeleteSystemRole) {
			helper.WriteErrorResponse(c, helper.NewBadRequestError("role", "Cannot delete system role"))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// GetPaginatedRoles godoc
// @Summary Get paginated roles
// @Tags roles
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(0)
// @Param pageSize query int false "Page size" default(10)
// @Param sort query string false "Sort field"
// @Param order query string false "Sort order" Enums(asc, desc) default(desc)
// @Success 200 {object} helper.AppSuccess{data=helper.PaginatedResponse{items=[]dto.RoleResponse, totalCount=int}}
// @Failure 400 {object} helper.AppError "Invalid input"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Router /admin/roles [get]
// @Security BearerAuth
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
	// Convert pointer slice to value slice for the new mapper
	roleValues := make([]domain.Role, len(domainRoles))
	for i, rolePtr := range domainRoles {
		roleValues[i] = *rolePtr
	}

	response := helper.PaginatedResponse{
		Items:      h.RoleMapper.DomainListToDTO(roleValues),
		TotalCount: total,
	}

	helper.WriteSuccessResponse(c, http.StatusOK, response, "Roles retrieved successfully")
}
