package handler

import (
	"net/http"
	"strconv"
	"time"

	user "github.com/EdwinRincon/browersfc-api/api/model"
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
	var err error
	pageStr := c.DefaultQuery("page", "0")

	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		page = 0
	}

	userRepo, err := h.UserRepository.ListUsers(c, page)
	if err != nil {
		helper.HandleError(c, http.StatusInternalServerError, "error al obtener usuarios", err)
		return
	}

	helper.ResponseJSONSuccess(c, "success", userRepo)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user user.Users

	if err := c.ShouldBindJSON(&user); err != nil {
		helper.HandleError(c, http.StatusBadRequest, "datos erróneos para crear el usuario", err)
		return
	}

	// Format Birthdate from YYYY-MM-DD to RFC3339
	birthdate := user.Birthdate.Format("2006-01-02T15:04:05Z")
	user.Birthdate, _ = time.Parse(time.RFC3339, birthdate)

	ok := h.UserRepository.CreateUser(c, &user)
	if !ok {
		return
	}
	helper.ResponseJSONSuccess(c, "success", user)
}
