package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
	"github.com/EdwinRincon/browersfc-api/api/service"
	"github.com/EdwinRincon/browersfc-api/config"
	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
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
		helper.HandleValidationError(c, err)
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
// @Param user body model.User true "User details"
// @Success 201 {object} model.UserMin "User created successfully"
// @Failure 400 {object} helper.AppError "Invalid input provided"
// @Failure 500 {object} helper.AppError "Internal server error occurred"
// @Router /users [post]
// @Security Bearer
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		helper.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	createdUser, err := h.UserService.CreateUser(ctx, &user)
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusCreated, createdUser, "User created successfully")
}

func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	ctx := c.Request.Context()
	user, err := h.UserService.GetUserByUsername(ctx, username)
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}
	if user == nil {
		helper.HandleError(c, helper.NewAppError(http.StatusNotFound, constants.ErrUserNotFound.Error(), ""), false)
		return
	}

	userResponse := model.UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		LastName:   user.LastName,
		Username:   user.Username,
		IsActive:   user.IsActive,
		Birthdate:  user.Birthdate,
		ImgProfile: user.ImgProfile,
		ImgBanner:  user.ImgBanner,
		RoleName:   user.Role.Name,
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
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, users, "Users listed successfully")
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var userUpdate model.UserUpdate
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		helper.HandleValidationError(c, err)
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
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusNoContent, nil, "User deleted successfully")
}

// LoginWithGoogle initiates OAuth2 flow with Google
func (h *UserHandler) LoginWithGoogle(c *gin.Context) {
	// Generate random state
	state, err := helper.GenerateRandomState()
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Authentication initialization failed", err.Error()), false)
		return
	}

	// Generate PKCE parameters
	pkceParams, err := config.GeneratePKCE()
	if err != nil {
		logrus.WithError(err).Error("Failed to generate PKCE parameters")
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Authentication initialization failed", err.Error()), false)
		return
	}

	// Store PKCE parameters
	config.StorePKCE(state, pkceParams)

	// Set state cookie with secure parameters
	c.SetCookie("oauth_state", state, 600, "/", "", true, true) // 10 minutes expiry

	// Generate authorization URL with PKCE
	opts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("code_challenge", pkceParams.Challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	}
	url := config.OAuthConfig.AuthCodeURL(state, opts...)

	helper.HandleSuccess(c, http.StatusOK, gin.H{"url": url}, "OAuth URL generated")
}

type GoogleUserInfo struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
	FamilyName string `json:"family_name"`
}

// GoogleCallback handles the OAuth2 callback
func (h *UserHandler) GoogleCallback(c *gin.Context) {
	// Validate OAuth state
	if !h.validateOAuthState(c) {
		return
	}

	// Get token from OAuth exchange
	token, err := h.exchangeCodeForToken(c)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusUnauthorized, "Authentication failed", err.Error()), false)
		return
	}

	// Get Google user info
	googleUser, err := h.fetchGoogleUserInfo(c.Request.Context(), token)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusUnauthorized, "Failed to get user info", err.Error()), false)
		return
	}

	// Check if user exists or create new one
	ctx := c.Request.Context()
	user, err := h.UserService.GetUserByUsername(ctx, googleUser.Email)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to check user", err.Error()), false)
		return
	}

	if user == nil {
		// Generate random password for new user
		randomPass, err := helper.GenerateRandomPassword()
		if err != nil {
			helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to create user", err.Error()), false)
			return
		}

		// Create new user
		newUser := &model.User{
			Username:   googleUser.Email,
			Name:       googleUser.Name,
			ImgProfile: googleUser.Picture,
			LastName:   googleUser.FamilyName,
			IsActive:   true,
			RoleID:     4,
			Password:   randomPass,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		var createdUserMin *model.UserMin
		createdUserMin, err = h.UserService.CreateUser(ctx, newUser)
		if err != nil {
			helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to create user", err.Error()), false)
			return
		}

		user, err = h.UserService.GetUserByUsername(ctx, createdUserMin.Username)
		if err != nil {
			helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to retrieve user", err.Error()), false)
			return
		}
	}

	if user == nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to create or retrieve user", ""), false)
		return
	}

	// Now proceed with authentication
	h.setAuthenticationResponse(c, user)
}

func (h *UserHandler) validateOAuthState(c *gin.Context) bool {
	state := c.Query("state")
	storedState, _ := c.Cookie("oauth_state")

	if state == "" || state != storedState {
		logrus.Error("Invalid OAuth state")
		helper.HandleError(c, helper.NewAppError(http.StatusUnauthorized, "Invalid OAuth state", ""), false)
		return false
	}

	pkceParams, ok := config.GetAndDeletePKCE(state)
	if !ok {
		logrus.Error("PKCE parameters not found")
		helper.HandleError(c, helper.NewAppError(http.StatusUnauthorized, "Invalid request", ""), false)
		return false
	}

	c.Set("pkce_params", pkceParams)
	return true
}

func (h *UserHandler) exchangeCodeForToken(c *gin.Context) (*oauth2.Token, error) {
	code := c.Query("code")
	pkceParams := c.MustGet("pkce_params").(*config.PKCEParams)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	token, err := config.OAuthConfig.Exchange(ctx, code,
		oauth2.SetAuthURLParam("code_verifier", pkceParams.Verifier))
	if err != nil {
		logrus.WithError(err).Error("Code exchange failed")
		return nil, err
	}

	if !token.Valid() {
		return nil, errors.New("invalid token received")
	}

	return token, nil
}

func (h *UserHandler) fetchGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := config.OAuthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var googleUser GoogleUserInfo
	if err := json.Unmarshal(userData, &googleUser); err != nil {
		return nil, err
	}

	return &googleUser, nil
}

func (h *UserHandler) setAuthenticationResponse(c *gin.Context, user *model.User) {
	jwtToken, err := h.AuthService.GenerateToken(user.Username, user.Role.Name)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to generate token", err.Error()), false)
		return
	}

	userResponse := model.UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		LastName:   user.LastName,
		Username:   user.Username,
		IsActive:   user.IsActive,
		ImgProfile: user.ImgProfile,
		RoleName:   user.Role.Name,
	}

	c.Header("Authorization", "Bearer "+jwtToken)
	c.SetCookie("token", jwtToken, 3600, "/", "", true, true)

	helper.HandleSuccess(c, http.StatusOK, gin.H{
		"token": jwtToken,
		"user":  userResponse,
	}, "User logged in successfully")
}
