package handler

import (
	"net/http"
	"strconv"

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
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid input", err.Error()), true)
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

	helper.HandleSuccess(c, http.StatusOK, gin.H{"token": token}, "User logged in successfully")
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user model.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid input", err.Error()), true)
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
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to retrieve user", err.Error()), false)
		return
	}
	if user == nil {
		helper.HandleError(c, helper.NewAppError(http.StatusNotFound, "User not found", ""), false)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, user, "User retrieved successfully")
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
	var user model.Users // Modelo completo del usuario con todos los campos
	var updatedUser *model.UserMin
	var err error

	// Bind JSON input to the user model
	if err = c.ShouldBindJSON(&user); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid input", err.Error()), true)
		return
	}

	// Use UserService to update the user
	ctx := c.Request.Context()
	if updatedUser, err = h.UserService.UpdateUser(ctx, &user); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to update user", err.Error()), true)
		return
	}

	// Check if updatedUser is nil, which might indicate the user was not found
	if updatedUser == nil {
		helper.HandleError(c, helper.NewAppError(http.StatusNotFound, "User not found", ""), true)
		return
	}

	// Respond with success and the updated user information
	helper.HandleSuccess(c, http.StatusOK, updatedUser, "User updated successfully")
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// validate uuid parse
	if _, err := uuid.Parse(id); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid input: ID must be valid", err.Error()), true)
		return
	}

	// Use UserService to delete the user
	ctx := c.Request.Context()
	err := h.UserService.DeleteUser(ctx, id)
	if err != nil {
		if err == repository.ErrUserNotFound {
			helper.HandleError(c, helper.NewAppError(http.StatusNotFound, "User not found", err.Error()), true)
		} else {
			helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to delete user", err.Error()), true)
		}
		return
	}

	helper.HandleSuccess(c, http.StatusOK, nil, "User deleted successfully")
}
