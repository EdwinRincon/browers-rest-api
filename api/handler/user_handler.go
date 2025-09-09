package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/EdwinRincon/browersfc-api/adapter/mapper"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/pkg/logger"
	"github.com/EdwinRincon/browersfc-api/pkg/security"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/config"
	"github.com/EdwinRincon/browersfc-api/helper"
	domainservice "github.com/EdwinRincon/browersfc-api/internal/domain/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	AuthenticationDomainService *domainservice.AuthenticationDomainService
	UserDomainService           *domainservice.UserDomainService
	RoleDomainService           *domainservice.RoleDomainService
	UserMapper                  *mapper.UserMapper
	googleClient                *http.Client // Pre-initialized Google API client for OAuth.
}

// GoogleUserInfo represents the user information returned by the Google API.
type GoogleUserInfo struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
	FamilyName string `json:"family_name"`
}

// NewUserHandler creates a new UserHandler with a pre-configured HTTP client.
func NewUserHandler(authService *domainservice.AuthenticationDomainService, userDomainService *domainservice.UserDomainService, roleDomainService *domainservice.RoleDomainService) *UserHandler {
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
		AuthenticationDomainService: authService,
		UserDomainService:           userDomainService,
		RoleDomainService:           roleDomainService,
		UserMapper:                  mapper.NewUserMapper(),
		googleClient:                client,
	}
}

// setAuthenticationResponse generates a JWT token and sets it in the response headers and cookies.
func (h *UserHandler) setAuthenticationResponse(c *gin.Context, user *domain.User) {
	// Get role name for JWT token
	roleName := ""
	if user.RoleID > 0 {
		role, err := h.RoleDomainService.GetRoleByID(c.Request.Context(), user.RoleID)
		if err == nil && role != nil {
			roleName = role.Name
		}
	}

	jwtToken, err := h.AuthenticationDomainService.GenerateToken(c.Request.Context(), user)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
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
	security.SetSecureCookie(c, "token", jwtToken, int(time.Hour/time.Second))

	helper.WriteSuccessResponse(c, http.StatusOK, gin.H{
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
		helper.WriteErrorResponse(c, helper.NewUnauthorizedError("Invalid OAuth state"))
		return false
	}

	pkceParams, ok := config.GetAndDeletePKCE(state)
	if !ok {
		helper.WriteErrorResponse(c, helper.NewUnauthorizedError("Invalid request"))
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
// @Summary Google OAuth2 callback
// @Tags users
// @ID googleCallback
// @Produce json
// @Success 200 {object} dto.UserResponse "Successful authentication and user data"
// @Failure 401 {object} helper.AppError "Authentication failed or invalid state"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Router /users/auth/google/callback [get]
func (h *UserHandler) GoogleCallback(c *gin.Context) {
	googleUser, err := h.performGoogleOAuth(c)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewUnauthorizedError("Authentication failed"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	user, err := h.UserDomainService.GetUserByUsername(ctx, googleUser.Email)
	if err != nil && !errors.Is(err, constants.ErrRecordNotFound) {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	if user == nil {
		// Validate email domain
		if !security.ValidateEmailDomain(googleUser.Email) {
			helper.WriteErrorResponse(c, helper.NewForbiddenError("Email domain not allowed"))
			return
		}

		// Note: Rate limiting for new accounts is now handled by middleware.RateLimitNewAccounts

		// Get the default role for new Google users
		defaultRole, err := h.RoleDomainService.GetRoleByName(ctx, constants.RoleDefault)
		if err != nil {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
			return
		}

		newUser := &domain.User{
			Username:   googleUser.Email,
			Name:       googleUser.Name,
			ImgProfile: googleUser.Picture,
			LastName:   googleUser.FamilyName,
			RoleID:     defaultRole.ID,
		}

		_, err = h.UserDomainService.CreateUser(ctx, newUser)
		if err != nil {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
			return
		}

		logger.Info(c, "new user created via OAuth", "username", newUser.Username)

		user, err = h.UserDomainService.GetUserByUsername(ctx, newUser.Username)
		if err != nil {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
			return
		}
	}

	h.setAuthenticationResponse(c, user)
}

// LoginWithGoogle godoc
// @Summary Initiate Google OAuth2 login
// @Tags users
// @ID loginWithGoogle
// @Produce json
// @Success 200 {object} map[string]string "Authorization URL generated"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Router /users/auth/google [get]
func (h *UserHandler) LoginWithGoogle(c *gin.Context) {
	state, err := helper.GenerateRandomState()
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	pkceParams, err := config.GeneratePKCE()
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	config.StorePKCE(state, pkceParams)

	security.SetSecureCookie(c, "oauth_state", state, int(10*time.Minute/time.Second))
	opts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("code_challenge", pkceParams.Challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	}
	url := config.OAuthConfig.AuthCodeURL(state, opts...)

	helper.WriteSuccessResponse(c, http.StatusOK, gin.H{"url": url}, "OAuth URL generated")
}

// CreateUser godoc
// @Summary Create a new user
// @Tags users
// @ID createUser
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User data"
// @Success 201 {object} dto.UserShort "User created successfully"
// @Failure 400 {object} helper.AppError "Invalid input"
// @Failure 409 {object} helper.AppError "Conflict (e.g., username exists)"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Router /admin/users [post]
// @Security BearerAuth
func (h *UserHandler) CreateUser(c *gin.Context) {
	var createRequest dto.CreateUserRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		helper.WriteErrorResponse(c, helper.BuildValidationErrorFromBinding(err, "body", constants.MsgInvalidUserData))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var roleID uint64
	if createRequest.RoleID > 0 {
		// If role_id is provided, verify it exists
		role, err := h.RoleDomainService.GetRoleByID(ctx, createRequest.RoleID)
		if err != nil {
			helper.WriteErrorResponse(c, helper.NewBadRequestError("role_id", "Role not found"))
			return
		}
		roleID = role.ID
	} else {
		// Get the default role if no role_id provided
		defaultRole, err := h.RoleDomainService.GetRoleByName(ctx, constants.RoleDefault)
		if err != nil {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
			return
		}
		roleID = defaultRole.ID
	}

	// Map the request to a domain User
	domainUser := h.UserMapper.DTOToDomain(&createRequest, roleID)

	createdUser, err := h.UserDomainService.CreateUser(ctx, domainUser)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrRecordAlreadyExists):
			helper.WriteErrorResponse(c, helper.NewConflictError("username", "Username already exists"))
			return
		default:
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
			return
		}
	}

	// Map domain user to DTO response
	response := h.UserMapper.DomainToShortDTO(createdUser)
	helper.WriteSuccessResponse(c, http.StatusCreated, response, "User created successfully")
}

// GetUserByUsername godoc
// @Summary Get a user by username
// @Tags users
// @ID getUserByUsername
// @Param username path string true "Username"
// @Success 200 {object} dto.UserResponse "User retrieved successfully"
// @Failure 400 {object} helper.AppError "Invalid input"
// @Failure 404 {object} helper.AppError "User not found"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Router /users/{username} [get]
// @Security BearerAuth
func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	user, err := h.UserDomainService.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("user"))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	// Map domain user to DTO response
	userResponse := h.UserMapper.DomainToDTO(user, nil)
	helper.WriteSuccessResponse(c, http.StatusOK, userResponse, "User retrieved successfully")
}

// GetPaginatedUsers godoc
// @Summary Get paginated users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(0)
// @Param pageSize query int false "Page size" default(10)
// @Param sort query string false "Sort field"
// @Param order query string false "Sort order" Enums(asc, desc) default(desc)
// @Success 200 {object} helper.AppSuccess{data=helper.PaginatedResponse{items=[]dto.UserResponse, totalCount=int}}
// @Failure 400 {object} helper.AppError "Invalid input"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Router /users [get]
// @Security BearerAuth
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
		helper.WriteErrorResponse(c, helper.NewBadRequestError("sort", err.Error()))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	users, total, err := h.UserDomainService.GetPaginatedUsers(ctx, sort, order, page, pageSize)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	response := helper.PaginatedResponse{
		Items:      h.UserMapper.DomainListToDTO(users),
		TotalCount: total,
	}

	helper.WriteSuccessResponse(c, http.StatusOK, response, "Users retrieved successfully")
}

// UpdateUser godoc
// @Summary Update an existing user
// @Tags users
// @ID updateUser
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param user body dto.UpdateUserRequest true "Updated user data"
// @Success 200 {object} dto.UserShort "User updated successfully"
// @Failure 400 {object} helper.AppError "Invalid input or UUID format"
// @Failure 404 {object} helper.AppError "User not found"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Router /admin/users/{id} [put]
// @Security BearerAuth
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	if _, err := uuid.Parse(userIDStr); err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid user ID format"))
		return
	}

	var userUpdateDTO dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&userUpdateDTO); err != nil {
		helper.WriteErrorResponse(c, helper.BuildValidationErrorFromBinding(err, "body", constants.MsgInvalidUserData))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Map DTO to domain entity
	domainUpdates := h.UserMapper.UpdateDTOToDomain(&userUpdateDTO)

	updatedUser, err := h.UserDomainService.UpdateUser(ctx, userIDStr, domainUpdates)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("user"))
			return
		}
		if errors.Is(err, constants.ErrRecordAlreadyExists) {
			helper.WriteErrorResponse(c, helper.NewConflictError("username", "Username already exists"))
			return
		}
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	// Map domain user to DTO response
	response := h.UserMapper.DomainToShortDTO(updatedUser)
	helper.WriteSuccessResponse(c, http.StatusOK, response, "User updated successfully")
}

// DeleteUser godoc
// @Summary Delete a user
// @Tags users
// @ID deleteUser
// @Param id path string true "User ID (UUID)"
// @Success 204 "No Content"
// @Failure 400 {object} helper.AppError "Invalid UUID format"
// @Router /admin/users/{id} [delete]
// @Security BearerAuth
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid UUID format"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.UserDomainService.DeleteUser(ctx, id.String())
	if err != nil && !errors.Is(err, constants.ErrRecordNotFound) {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	c.Status(http.StatusNoContent)
}
