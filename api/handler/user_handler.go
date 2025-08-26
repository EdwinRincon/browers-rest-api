package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/EdwinRincon/browersfc-api/api/auth"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/mapper"
	"github.com/EdwinRincon/browersfc-api/pkg/logger"
	"gorm.io/gorm"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/service"
	"github.com/EdwinRincon/browersfc-api/config"
	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	AuthService  service.AuthService
	UserService  service.UserService
	RoleService  service.RoleService
	googleClient *http.Client // Pre-initialized Google API client for OAuth.
}

// GoogleUserInfo represents the user information returned by the Google API.
type GoogleUserInfo struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
	FamilyName string `json:"family_name"`
}

// NewUserHandler creates a new UserHandler with a pre-configured HTTP client.
func NewUserHandler(authService service.AuthService, userService service.UserService, roleService service.RoleService) *UserHandler {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		},
	}

	return &UserHandler{
		AuthService:  authService,
		UserService:  userService,
		RoleService:  roleService,
		googleClient: client,
	}
}

// setAuthenticationResponse generates a JWT token and sets it in the response headers and cookies.
func (h *UserHandler) setAuthenticationResponse(c *gin.Context, user *model.User) {
	roleName := ""
	if user.Role != nil {
		roleName = user.Role.Name
	}
	jwtToken, err := h.AuthService.GenerateToken(user.Username, roleName)
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	authUserResponse := dto.AuthUserResponse{
		ID:         user.ID,
		Name:       user.Name,
		LastName:   user.LastName,
		Username:   user.Username,
		ImgProfile: user.ImgProfile,
		RoleName:   roleName,
	}

	c.Header("Authorization", "Bearer "+jwtToken)
	auth.SetSecureCookie(c, "token", jwtToken, int(time.Hour/time.Second))

	helper.HandleSuccess(c, http.StatusOK, gin.H{
		"token": jwtToken,
		"user":  authUserResponse,
	}, "User logged in successfully")
}

// fetchGoogleUserInfo retrieves user information from Google using the provided OAuth2 token.
func (h *UserHandler) fetchGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, err
	}
	token.SetAuthHeader(req)

	resp, err := h.googleClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Google API error: " + resp.Status)
	}

	var googleUser GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, err
	}

	return &googleUser, nil
}

// exchangeCodeForToken exchanges the authorization code for an OAuth2 token.
func (h *UserHandler) exchangeCodeForToken(c *gin.Context) (*oauth2.Token, error) {
	code := c.Query("code")
	pkceParams := c.MustGet("pkce_params").(*config.PKCEParams)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	token, err := config.OAuthConfig.Exchange(ctx, code,
		oauth2.SetAuthURLParam("code_verifier", pkceParams.Verifier))
	if err != nil {
		logger.Error(c, "code exchange failed", "error", err)
		return nil, err
	}

	if !token.Valid() {
		return nil, errors.New("invalid token received")
	}

	return token, nil
}

// validateOAuthState validates the OAuth state parameter and retrieves PKCE parameters.
func (h *UserHandler) validateOAuthState(c *gin.Context) bool {
	state := c.Query("state")
	storedState, _ := c.Cookie("oauth_state")

	logger.Info(c, "OAuth state validation", "query_state", state, "cookie_state", storedState)

	if state == "" || state != storedState {
		helper.RespondWithError(c, helper.Unauthorized("Invalid OAuth state"))
		return false
	}

	pkceParams, ok := config.GetAndDeletePKCE(state)
	if !ok {
		helper.RespondWithError(c, helper.Unauthorized("Invalid request"))
		return false
	}

	// Clear the state cookie immediately after successful validation
	//TODO: modify the cookie to be secure and HttpOnly when in production
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)
	c.Set("pkce_params", pkceParams)
	return true
}

// performGoogleOAuth orchestrates the OAuth flow: state validation, token exchange, and user info retrieval.
func (h *UserHandler) performGoogleOAuth(c *gin.Context) (*GoogleUserInfo, error) {
	if !h.validateOAuthState(c) {
		return nil, errors.New("invalid OAuth state")
	}

	token, err := h.exchangeCodeForToken(c)
	if err != nil {
		return nil, err
	}

	return h.fetchGoogleUserInfo(c.Request.Context(), token)
}

// GoogleCallback godoc
// @Summary      Google OAuth2 callback
// @Description  Handles the OAuth2 callback from Google and logs in or registers the user
// @Tags         users
// @ID           googleCallback
// @Produce      json
// @Success      200  {object}  dto.UserResponse "Successful authentication and user data"
// @Failure      401  {object}  helper.AppError "Authentication failed or invalid state"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /users/auth/google/callback [get]
func (h *UserHandler) GoogleCallback(c *gin.Context) {
	googleUser, err := h.performGoogleOAuth(c)
	if err != nil {
		helper.RespondWithError(c, helper.Unauthorized("Authentication failed"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	user, err := h.UserService.GetUserByUsername(ctx, googleUser.Email)
	if err != nil && !errors.Is(err, constants.ErrRecordNotFound) {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	if user == nil {
		// Validate email domain
		if !auth.ValidateEmailDomain(googleUser.Email) {
			helper.RespondWithError(c, helper.StatusForbidden("Email domain not allowed"))
			return
		}

		// Note: Rate limiting for new accounts is now handled by middleware.RateLimitNewAccounts

		// Get the default role for new Google users
		defaultRole, err := h.RoleService.GetActiveRoleByName(ctx, constants.RoleDefault)
		if err != nil {
			helper.RespondWithError(c, helper.InternalError(err))
			return
		}

		newUser := &model.User{
			Username:   googleUser.Email,
			Name:       googleUser.Name,
			ImgProfile: googleUser.Picture,
			LastName:   googleUser.FamilyName,
			RoleID:     defaultRole.ID,
		}

		_, err = h.UserService.CreateUser(ctx, newUser)
		if err != nil {
			helper.RespondWithError(c, helper.InternalError(err))
			return
		}

		logger.Info(c, "new user created via OAuth", "username", newUser.Username)

		user, err = h.UserService.GetUserByUsername(ctx, newUser.Username)
		if err != nil {
			helper.RespondWithError(c, helper.InternalError(err))
			return
		}
	}

	h.setAuthenticationResponse(c, user)
}

// LoginWithGoogle godoc
// @Summary      Initiate Google OAuth2 login
// @Description  Initiates the OAuth2 flow with Google and returns the authorization URL
// @Tags         users
// @ID           loginWithGoogle
// @Produce      json
// @Success      200  {object}  map[string]string "Authorization URL generated"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /users/auth/google [get]
func (h *UserHandler) LoginWithGoogle(c *gin.Context) {
	state, err := helper.GenerateRandomState()
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	pkceParams, err := config.GeneratePKCE()
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	config.StorePKCE(state, pkceParams)

	auth.SetSecureCookie(c, "oauth_state", state, int(10*time.Minute/time.Second))
	opts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("code_challenge", pkceParams.Challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	}
	url := config.OAuthConfig.AuthCodeURL(state, opts...)

	helper.HandleSuccess(c, http.StatusOK, gin.H{"url": url}, "OAuth URL generated")
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Creates a new user with the provided data. If role_id is not provided, assigns the default user role.
// @Tags         users
// @ID           createUser
// @Accept       json
// @Produce      json
// @Param        user  body      dto.CreateUserRequest  true  "User data"
// @Success      201   {object}  dto.UserShort  "User created successfully"
// @Failure      400   {object}  helper.AppError "Invalid input"
// @Failure      409   {object}  helper.AppError "Conflict (e.g., username exists)"
// @Failure      500   {object}  helper.AppError "Internal server error"
// @Router       /users [post]
// @Security     ApiKeyAuth
func (h *UserHandler) CreateUser(c *gin.Context) {
	var createRequest dto.CreateUserRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		helper.RespondWithError(c, helper.ProcessValidationError(err, "body", constants.MsgInvalidUserData))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var roleID uint8
	if createRequest.RoleID > 0 {
		// If role_id is provided, verify it exists
		role, err := h.RoleService.GetRoleByID(ctx, createRequest.RoleID)
		if err != nil {
			helper.RespondWithError(c, helper.BadRequest("role_id", "Role not found"))
			return
		}
		roleID = role.ID
	} else {
		// Get the default role if no role_id provided
		defaultRole, err := h.RoleService.GetActiveRoleByName(ctx, constants.RoleDefault)
		if err != nil {
			helper.RespondWithError(c, helper.InternalError(err))
			return
		}
		roleID = defaultRole.ID
	}

	// Map the request to a User model
	user := mapper.ToUser(&createRequest, roleID)

	createdUser, err := h.UserService.CreateUser(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrRecordAlreadyExists):
			helper.RespondWithError(c, helper.Conflict("username", "Username already exists"))
			return
		default:
			helper.RespondWithError(c, helper.InternalError(err))
			return
		}
	}

	helper.HandleSuccess(c, http.StatusCreated, createdUser, "User created successfully")
}

// GetUserByUsername godoc
// @Summary      Get a user by username
// @Description  Retrieves a user by their username
// @Tags         users
// @ID           getUserByUsername
// @Param        username  path      string  true  "Username"
// @Success      200      {object}  dto.UserResponse "User retrieved successfully"
// @Failure      400      {object}  helper.AppError "Invalid input"
// @Failure      404      {object}  helper.AppError "User not found"
// @Failure      500      {object}  helper.AppError "Internal server error"
// @Router       /users/{username} [get]
// @Security     ApiKeyAuth
func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	user, err := h.UserService.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("user"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}
	if user == nil {
		helper.RespondWithError(c, helper.NotFound("user"))
		return
	}

	userResponse := mapper.ToUserResponse(user)

	helper.HandleSuccess(c, http.StatusOK, userResponse, "User retrieved successfully")
}

// GetPaginatedUsers godoc
// @Summary      Get paginated users
// @Description  Retrieves a paginated list of users with sorting and ordering
// @Tags         users
// @ID           getPaginatedUsers
// @Param        page      query     int     false  "Page number (0-based)"
// @Param        pageSize  query     int     false  "Number of items per page (default 10)"
// @Param        sort      query     string  false  "Sort field (e.g., username, name)"
// @Param        order     query     string  false  "Sort order (asc/desc)"
// @Success      200       {object}  map[string]interface{} "Users retrieved successfully"
// @Failure      400       {object}  helper.AppError "Invalid input"
// @Failure      500       {object}  helper.AppError "Internal server error"
// @Router       /users [get]
// @Security     ApiKeyAuth
func (h *UserHandler) GetPaginatedUsers(c *gin.Context) {
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 0 {
		page = 0
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 10
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	// Validate sort field
	if err := helper.ValidateSort(model.User{}, sort); err != nil {
		helper.RespondWithError(c, helper.BadRequest("sort", err.Error()))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	users, total, err := h.UserService.GetPaginatedUsers(ctx, sort, order, page, pageSize)
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	response := helper.PaginatedResponse{
		Items:      mapper.ToUserResponseList(users),
		TotalCount: total,
	}

	helper.HandleSuccess(c, http.StatusOK, response, "Users retrieved successfully")
}

// UpdateUser godoc
// @Summary      Update an existing user
// @Description  Updates an existing user's information by ID
// @Tags         users
// @ID           updateUser
// @Accept       json
// @Produce      json
// @Param        id    path      string           true  "User ID (UUID)"
// @Param        user  body      dto.UpdateUserRequest true  "Updated user data"
// @Success      200   {object}  dto.UserShort  "User updated successfully"
// @Failure      400   {object}  helper.AppError "Invalid input or UUID format"
// @Failure      404   {object}  helper.AppError "User not found"
// @Failure      500   {object}  helper.AppError "Internal server error"
// @Router       /users/{id} [put]
// @Security     ApiKeyAuth
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	if _, err := uuid.Parse(userIDStr); err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid user ID format"))
		return
	}

	var userUpdateDTO dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&userUpdateDTO); err != nil {
		helper.RespondWithError(c, helper.ProcessValidationError(err, "body", constants.MsgInvalidUserData))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	updatedUser, err := h.UserService.UpdateUser(ctx, &userUpdateDTO, userIDStr)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("user"))
			return
		}
		if errors.Is(err, constants.ErrRecordAlreadyExists) {
			helper.RespondWithError(c, helper.Conflict("username", "Username already exists"))
			return
		}
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	response := mapper.ToUserShort(updatedUser)
	helper.HandleSuccess(c, http.StatusOK, response, "User updated successfully")
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Deletes a user by their ID
// @Tags         users
// @ID           deleteUser
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      204 "No Content"
// @Failure      400  {object}  helper.AppError "Invalid UUID format"
// @Router       /users/{id} [delete]
// @Security     ApiKeyAuth
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid UUID format"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	err = h.UserService.DeleteUser(ctx, id.String())
	if err != nil && !errors.Is(err, constants.ErrRecordNotFound) {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	c.Status(http.StatusNoContent)
}
