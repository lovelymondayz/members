package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lovelymondayz/members/backend/src/models"
	"github.com/lovelymondayz/members/backend/src/repository"
)

type UpdateMemberRequest struct {
	MemberCode string `json:"member_code"`
	Tier       string `json:"tier"`
}

type MemberHandler struct {
	repo *repository.MemberRepository
}

func NewMemberHandler(repo *repository.MemberRepository) *MemberHandler {
	return &MemberHandler{repo: repo}
}

// GetMembers godoc
// @Summary List members for a store
// @Tags members
// @Param store_id query string true "Store ID"
// @Success 200 {array} models.Member
// @Router /members [get]
// @Security BearerAuth
func (h *MemberHandler) GetMembers(c *gin.Context) {
	// Enforce store_id from JWT (data isolation)
	storeID, hasStore := c.Get("store_id")
	if hasStore {
		members, err := h.repo.FindByStoreID(storeID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch members"})
			return
		}
		c.JSON(http.StatusOK, members)
		return
	}

	// Super admin can query any store
	queryStoreID := c.Query("store_id")
	if queryStoreID != "" {
		members, err := h.repo.FindByStoreID(queryStoreID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch members"})
			return
		}
		c.JSON(http.StatusOK, members)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "store_id is required"})
}

// CreateMember godoc
// @Summary Create a new member
// @Tags members
// @Param member body models.Member true "Member data"
// @Success 201 {object} models.Member
// @Router /members [post]
// @Security BearerAuth
func (h *MemberHandler) CreateMember(c *gin.Context) {
	var member models.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Enforce store_id from JWT (data isolation)
	storeID, hasStore := c.Get("store_id")
	if hasStore {
		parsedID, err := uuid.Parse(storeID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid store_id in token"})
			return
		}
		member.StoreID = parsedID
	} else {
		// Super admin must provide store_id in body
		if member.StoreID == uuid.Nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "store_id is required"})
			return
		}
	}

	if err := h.repo.Create(&member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create member"})
		return
	}

	c.JSON(http.StatusCreated, member)
}

// GetMember godoc
// @Summary Get member by ID
// @Tags members
// @Param id path string true "Member ID"
// @Success 200 {object} models.Member
// @Router /members/{id} [get]
// @Security BearerAuth
func (h *MemberHandler) GetMember(c *gin.Context) {
	id := c.Param("id")

	member, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
		return
	}

	c.JSON(http.StatusOK, member)
}

// UpdateMember godoc
// @Summary Update member
// @Tags members
// @Param id path string true "Member ID"
// @Param request body UpdateMemberRequest true "Updated data"
// @Success 200 {object} models.Member
// @Router /members/{id} [put]
// @Security BearerAuth
func (h *MemberHandler) UpdateMember(c *gin.Context) {
	id := c.Param("id")

	var req UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
		return
	}

	if req.Tier != "" {
		member.Tier = req.Tier
	}
	if req.MemberCode != "" {
		member.MemberCode = req.MemberCode
	}

	if err := h.repo.Update(member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update member"})
		return
	}

	c.JSON(http.StatusOK, member)
}

// DeleteMember godoc
// @Summary Delete member (soft delete)
// @Tags members
// @Param id path string true "Member ID"
// @Success 204
// @Router /members/{id} [delete]
// @Security BearerAuth
func (h *MemberHandler) DeleteMember(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete member"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetMemberCard godoc
// @Summary Get membership card data (for digital card display)
// @Tags members
// @Param id path string true "Member ID"
// @Success 200 {object} map[string]interface{}
// @Router /members/{id}/card [get]
// @Security BearerAuth
func (h *MemberHandler) GetMemberCard(c *gin.Context) {
	id := c.Param("id")

	member, err := h.repo.FindByIDWithRelations(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
		return
	}

	cardData := map[string]interface{}{
		"member_id":    member.MemberID,
		"member_code":  member.MemberCode,
		"store_id":     member.StoreID,
		"tier":         member.Tier,
		"joined_at":    member.JoinedAt,
		"qr_data":      member.MemberCode,
		"name":         member.UserName,
		"store_name":   member.StoreName,
		"card_color":   member.StoreCardColor,
	}

	c.JSON(http.StatusOK, cardData)
}
