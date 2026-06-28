package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lovelymondayz/members/backend/src/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// Login godoc
// @Summary Login with email and password
// @Tags auth
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} AuthResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authService.LoginWithPassword(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User: map[string]interface{}{
			"id":    user.UserID,
			"name":  user.Name,
			"email": user.Email,
			"role":  roleName(user.RoleID),
		},
	})
}

// Register godoc
// @Summary Register new user
// @Tags auth
// @Param request body RegisterRequest true "Registration data"
// @Success 201 {object} AuthResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authService.RegisterUser(req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User: map[string]interface{}{
			"id":    user.UserID,
			"name":  user.Name,
			"email": user.Email,
			"role":  roleName(user.RoleID),
		},
	})
}

func roleName(roleID uint) string {
	switch roleID {
	case 1:
		return "super_admin"
	case 2:
		return "admin"
	case 3:
		return "member"
	default:
		return "unknown"
	}
}

// GoogleLogin godoc
// @Summary Start Google OAuth flow
// @Tags auth
// @Success 302
// @Router /auth/google [get]
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	url := h.authService.GetGoogleAuthURL()
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback godoc
// @Summary Google OAuth callback
// @Tags auth
// @Param query state true "OAuth state"
// @Param query code true "OAuth code"
// @Success 200 {object} AuthResponse
// @Router /auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}

	user, token, err := h.authService.HandleGoogleCallback(code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Redirect frontend with token (in production, use httpOnly cookie or postMessage)
	frontendURL := c.DefaultQuery("redirect", "http://localhost:3001")
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?token=%s", frontendURL, token))
	_ = user
}

// Me godoc
// @Summary Get current user profile
// @Tags auth
// @Success 200 {object} models.User
// @Router /auth/me [get]
// @Security BearerAuth
func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	storeID, _ := c.Get("store_id")

	c.JSON(http.StatusOK, gin.H{
		"user_id":  userID,
		"role":     role,
		"store_id": storeID,
	})
}

// Logout godoc
// @Summary Logout
// @Tags auth
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
