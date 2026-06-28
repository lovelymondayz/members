package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lovelymondayz/members/backend/src/config"
	"github.com/lovelymondayz/members/backend/src/models"
	"github.com/lovelymondayz/members/backend/src/service"
)

type AdminHandler struct {
	authService *service.AuthService
}

func NewAdminHandler(authService *service.AuthService) *AdminHandler {
	return &AdminHandler{authService: authService}
}

type CreateAdminRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// CreateAdmin godoc
// @Summary Create a new admin with store (Super Admin only)
// @Tags admin
// @Param request body CreateAdminRequest true "Admin data"
// @Success 201 {object} map[string]interface{}
// @Router /admin/admins [post]
// @Security BearerAuth
func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var req CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin, store, err := h.authService.CreateAdminWithStore(req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"admin": map[string]interface{}{
			"id":    admin.UserID,
			"name":  admin.Name,
			"email": admin.Email,
		},
		"store": map[string]interface{}{
			"id":   store.StoreID,
			"name": store.Name,
		},
	})
}

// GetAdmins godoc
// @Summary List all admins (Super Admin only)
// @Tags admins
// @Success 200 {array} models.User
// @Router /admin/admins [get]
// @Security BearerAuth
func (h *AdminHandler) GetAdmins(c *gin.Context) {
	var users []models.User
	err := config.DB.Where("role_id IN ?", []uint{1, 2}).Order("created_at DESC").Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch admins"})
		return
	}

	var result []gin.H
	for _, u := range users {
		result = append(result, gin.H{
			"user_id":   u.UserID,
			"name":      u.Name,
			"email":     u.Email,
			"role_id":   u.RoleID,
			"role":      roleName(u.RoleID),
			"is_active": u.IsActive,
		})
	}
	c.JSON(http.StatusOK, result)
}

// DashboardStats godoc
// @Summary Platform-wide statistics (Super Admin only)
// @Tags admin
// @Success 200 {object} map[string]interface{}
// @Router /admin/dashboard [get]
// @Security BearerAuth
func (h *AdminHandler) DashboardStats(c *gin.Context) {
	var totalStores int64
	var totalMembers int64
	var totalRevenue float64

	config.DB.Model(&models.Store{}).Count(&totalStores)
	config.DB.Model(&models.Member{}).Count(&totalMembers)
	config.DB.Model(&models.Payment{}).Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue)

	c.JSON(http.StatusOK, gin.H{
		"total_stores":  totalStores,
		"total_members": totalMembers,
		"total_revenue": totalRevenue,
	})
}
