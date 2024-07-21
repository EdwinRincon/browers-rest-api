package handler

import (
	"net/http"
	"strconv"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserRepository repository.UserRepository
}

func NewUserHandler(ur repository.UserRepository) *UserHandler {
	return &UserHandler{
		UserRepository: ur,
	}
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1") // Default to page 1 instead of 0

	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid page number", err.Error()))
		return
	}

	ctx := c.Request.Context()
	users, err := h.UserRepository.ListUsers(ctx, page)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to list users", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user model.Users
	var createdUser *model.UsersResponse
	var err error

	if err = c.ShouldBindJSON(&user); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid input", err.Error()))
		return
	}

	ctx := c.Request.Context()
	if createdUser, err = h.UserRepository.CreateUser(ctx, &user); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to create user", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    createdUser,
	})
}
