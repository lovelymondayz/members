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
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password,omitempty"`
	StoreName    string `json:"store_name,omitempty"`
	Address      string `json:"address,omitempty"`
	Phone        string `json:"phone,omitempty"`
	CardColorHex string `json:"card_color_hex,omitempty"`
}

type UpdateAdminRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
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

	admin, store, err := h.authService.CreateAdminWithStore(
		req.Name, req.Email, req.Password,
		req.StoreName, req.Address, req.Phone, req.CardColorHex,
	)
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
			"id":             store.StoreID,
			"name":           store.Name,
			"address":        store.Address,
			"phone":          store.Phone,
			"card_color_hex": store.CardColorHex,
		},
	})
}

// UpdateAdmin godoc
// @Summary Update an admin (Super Admin only)
// @Tags admin
// @Param id path string true "Admin User ID"
// @Param request body UpdateAdminRequest true "Admin data"
// @Success 200 {object} map[string]interface{}
// @Router /admin/admins/:id [put]
// @Security BearerAuth
func (h *AdminHandler) UpdateAdmin(c *gin.Context) {
	adminID := c.Param("id")

	var user models.User
	if err := config.DB.Where("user_id = ?", adminID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin not found"})
		return
	}

	var req UpdateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update admin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.UserID,
		"name":  user.Name,
		"email": user.Email,
		"role":  roleName(user.RoleID),
	})
}

// DeleteAdmin godoc
// @Delete an admin and their store (Super Admin only)
// @Tags admin
// @Param id path string true "Admin User ID"
// @Success 200 {object} map[string]string
// @Router /admin/admins/:id [delete]
// @Security BearerAuth
func (h *AdminHandler) DeleteAdmin(c *gin.Context) {
	adminID := c.Param("id")

	var user models.User
	if err := config.DB.Where("user_id = ?", adminID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin not found"})
		return
	}

	// Prevent self-deletion
	currentUserID, _ := c.Get("user_id")
	if currentUserID == user.UserID.String() {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete yourself"})
		return
	}

	// Delete store first (FK constraint)
	config.DB.Where("admin_id = ?", adminID).Delete(&models.Store{})

	// Delete the admin user
	config.DB.Delete(&user)

	c.JSON(http.StatusOK, gin.H{"message": "admin deleted"})
}

// GetAdmins godoc
// @Summary List all admins (Super Admin only)
// @Tags admins
// @Success 200 {array} models.User
// @Router /admin/admins [get]
// @Security BearerAuth
func (h *AdminHandler) GetAdmins(c *gin.Context) {
	var users []models.User
	err := config.DB.Where("role_id = ?", 2).Order("created_at DESC").Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch admins"})
		return
	}

	var result []gin.H
	for _, u := range users {
		// Find store for this admin
		var store models.Store
		config.DB.Where("admin_id = ?", u.UserID).First(&store)

		result = append(result, gin.H{
			"user_id":        u.UserID,
			"name":           u.Name,
			"email":          u.Email,
			"role_id":        u.RoleID,
			"role":           roleName(u.RoleID),
			"is_active":      u.IsActive,
			"store_id":       store.StoreID,
			"store_name":     store.Name,
			"store_address":  store.Address,
			"store_phone":    store.Phone,
			"store_color":    store.CardColorHex,
			"created_at":     u.CreatedAt,
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
