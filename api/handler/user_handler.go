package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
	"github.com/EdwinRincon/browersfc-api/api/service"
	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	AuthService service.AuthService
	UserService service.UserService
}

func NewUserHandler(authService service.AuthService, userService service.UserService) *UserHandler {
	return &UserHandler{
		AuthService: authService,
		UserService: userService,
	}
}

// LoginUser is a handler function to login a user
func (h *UserHandler) Login(c *gin.Context) {
	var user model.UserLogin
	if err := c.ShouldBindJSON(&user); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	// Verificar credenciales usando AuthService
	ctx := c.Request.Context()
	token, err := h.AuthService.Authenticate(ctx, user.Username, user.Password)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusUnauthorized, "Invalid username or password", err.Error()), false)
		return
	}

	// Configurar encabezado Authorization
	c.Header("Authorization", "Bearer "+token)

	// Configurar cookie
	c.SetCookie("token", token, 3600, "/", "", true, true)

	helper.HandleSuccess(c, http.StatusOK, gin.H{"token": token}, "User logged in successfully")
}

// CreateUser godoc
// @Summary Create a new user
// @Description This endpoint allows for the creation of a new user with the specified details.
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.Users true "User details"
// @Success 201 {object} model.UserMin "User created successfully"
// @Failure 400 {object} helper.AppError "Invalid input provided"
// @Failure 500 {object} helper.AppError "Internal server error occurred"
// @Router /users [post]
// @Security Bearer
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user model.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	createdUser, err := h.UserService.CreateUser(ctx, &user)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to create user", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusCreated, createdUser, "User created successfully")
}

func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	ctx := c.Request.Context()
	user, err := h.UserService.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			helper.HandleError(c, helper.NewAppError(http.StatusNotFound, constants.ErrUserNotFound.Error(), ""), false)
			return
		}
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to retrieve user", err.Error()), false)
		return
	}
	if user == nil {
		helper.HandleError(c, helper.NewAppError(http.StatusNotFound, constants.ErrUserNotFound.Error(), ""), false)
		return
	}

	userResponse := model.UsersResponse{
		ID:         user.ID,
		Name:       user.Name,
		LastName:   user.LastName,
		Username:   user.Username,
		IsActive:   user.IsActive,
		Birthdate:  user.Birthdate,
		ImgProfile: user.ImgProfile,
		ImgBanner:  user.ImgBanner,
		RoleName:   user.Roles.Name,
	}

	helper.HandleSuccess(c, http.StatusOK, userResponse, "User retrieved successfully")
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1") // Default to page 1 instead of 0

	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid page number", err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	users, err := h.UserService.ListUsers(ctx, page)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to list users", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, users, "Users listed successfully")
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var userUpdate model.UserUpdate
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	updatedUser, err := h.UserService.UpdateUser(ctx, &userUpdate, c.Param("id"))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			helper.HandleError(c, helper.NewAppError(http.StatusNotFound, constants.ErrUserNotFound.Error(), ""), true)
			return
		}
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to update user", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, updatedUser, "User updated successfully")
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid input: ID must be a valid UUID", err.Error()), true)
		return
	}

	// Attempt to delete the user
	ctx := c.Request.Context()
	err := h.UserService.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			helper.HandleError(c, helper.NewAppError(http.StatusNotFound, constants.ErrUserNotFound.Error(), ""), true)
			return
		}
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to delete user", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusNoContent, nil, "User deleted successfully")
}
