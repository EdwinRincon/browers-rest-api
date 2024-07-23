package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/EdwinRincon/browersfc-api/pkg/jwt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UserRepository repository.UserRepository
}

func NewUserHandler(ur repository.UserRepository) *UserHandler {
	return &UserHandler{
		UserRepository: ur,
	}
}

// LoginUser is a handler function to login a user
func (h *UserHandler) Login(c *gin.Context) {
	var user model.UserLogin
	var token string
	var err error

	if err = c.ShouldBindJSON(&user); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid input", err.Error()))
		return
	}

	// Verificar credenciales
	ctx := c.Request.Context()
	storedUser, err := h.UserRepository.GetUserByUsername(ctx, user.Username)
	if err != nil {
		// no usar helper.HandleError porque no queremos mostrar el error real al usuario
		c.JSON(http.StatusUnauthorized, gin.H{"code:": http.StatusUnauthorized, "error": "Invalid username or password"})
		return
	}

	// Verificar el numero de intentos fallidos de inicio de sesión (FailedLoginAttempts)
	if storedUser.FailedLoginAttempts >= 5 {
		// no usar helper.HandleError porque no queremos mostrar el error real al usuario
		c.JSON(http.StatusUnauthorized, gin.H{"code:": http.StatusUnauthorized, "error": "Invalid username or password"})
		return
	}

	// Comparar contraseñas
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		// Incrementar el número de intentos fallidos de inicio de sesión
		storedUser.FailedLoginAttempts++
		if _, err = h.UserRepository.UpdateUser(ctx, storedUser); err != nil {
			helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed login", err.Error()))
			return
		}
		// no usar helper.HandleError porque no queremos mostrar el error real al usuario
		c.JSON(http.StatusUnauthorized, gin.H{"code:": http.StatusUnauthorized, "error": "Invalid username or password"})
		return
	}

	// Restablecer el número de intentos fallidos de inicio de sesión
	storedUser.FailedLoginAttempts = 0
	if _, err = h.UserRepository.UpdateUser(ctx, storedUser); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed update login", err.Error()))
		return
	}

	// Generar token
	if token, err = jwt.GenerateToken(user.Username, storedUser.Roles.Name); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to generate token", err.Error()))
		return
	}

	// Configurar cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(2 * time.Hour),
		HttpOnly: true,
		Secure:   true, // Asegúrate de usar HTTPS en producción
		SameSite: http.SameSiteStrictMode,
	})

	helper.HandleSuccess(c, http.StatusOK, gin.H{"token": token}, "User logged in successfully")

}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user model.Users
	var createdUser *model.UserMin
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

	helper.HandleSuccess(c, http.StatusCreated, createdUser, "User created successfully")
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

	helper.HandleSuccess(c, http.StatusOK, users, "Users listed successfully")
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var user model.Users
	var updatedUser *model.UserMin
	var err error

	if err = c.ShouldBindJSON(&user); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid input", err.Error()))
		return
	}

	ctx := c.Request.Context()
	if updatedUser, err = h.UserRepository.UpdateUser(ctx, &user); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to update user", err.Error()))
		return
	}

	helper.HandleSuccess(c, http.StatusOK, updatedUser, "User updated successfully")
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	var user struct {
		Username string `json:"username" binding:"required"`
	}
	var err error

	if err = c.ShouldBindJSON(&user); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid input", err.Error()))
		return
	}

	ctx := c.Request.Context()
	if err = h.UserRepository.DeleteUser(ctx, user.Username); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to delete user", err.Error()))
		return
	}

	helper.HandleSuccess(c, http.StatusOK, nil, "User deleted successfully")
}
