package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lovelymondayz/members/backend/src/models"
	"github.com/lovelymondayz/members/backend/src/repository"
)

type StoreHandler struct {
	repo *repository.StoreRepository
}

func NewStoreHandler(repo *repository.StoreRepository) *StoreHandler {
	return &StoreHandler{repo: repo}
}

// GetStores godoc
// @Summary List all stores (Super Admin only)
// @Tags stores
// @Success 200 {array} models.Store
// @Router /stores [get]
// @Security BearerAuth
func (h *StoreHandler) GetStores(c *gin.Context) {
	stores, err := h.repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch stores"})
		return
	}
	c.JSON(http.StatusOK, stores)
}

// CreateStore godoc
// @Summary Create a new store (Super Admin only)
// @Tags stores
// @Param request body CreateStoreRequest true "Store data"
// @Success 201 {object} models.Store
// @Router /stores [post]
// @Security BearerAuth
func (h *StoreHandler) CreateStore(c *gin.Context) {
	var req CreateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store := models.Store{
		StoreID:      uuid.New(),
		AdminID:      uuid.MustParse(req.AdminID),
		Name:         req.Name,
		LogoURL:      req.LogoURL,
		Address:      req.Address,
		Phone:        req.Phone,
		CardColorHex: req.CardColorHex,
	}

	if store.CardColorHex == "" {
		store.CardColorHex = "#1E40AF"
	}

	if err := h.repo.Create(&store); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create store"})
		return
	}

	c.JSON(http.StatusCreated, store)
}

type CreateStoreRequest struct {
	AdminID      string `json:"admin_id" binding:"required,uuid"`
	Name         string `json:"name" binding:"required,min=2"`
	LogoURL      string `json:"logo_url,omitempty"`
	Address      string `json:"address,omitempty"`
	Phone        string `json:"phone,omitempty"`
	CardColorHex string `json:"card_color_hex,omitempty"`
}
