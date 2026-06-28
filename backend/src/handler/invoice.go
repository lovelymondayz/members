package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lovelymondayz/members/backend/src/models"
	"github.com/lovelymondayz/members/backend/src/repository"
)

type InvoiceHandler struct {
	invoiceRepo *repository.InvoiceRepository
	paymentRepo *repository.PaymentRepository
}

func NewInvoiceHandler(invoiceRepo *repository.InvoiceRepository, paymentRepo *repository.PaymentRepository) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceRepo: invoiceRepo,
		paymentRepo: paymentRepo,
	}
}

type CreateInvoiceRequest struct {
	MemberID    string  `json:"member_id" binding:"required"`
	StoreID     string  `json:"store_id"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
	DueDate     string  `json:"due_date"`
}

type RecordPaymentRequest struct {
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	Method    string  `json:"method"`
	Reference string  `json:"reference"`
	Note      string  `json:"note"`
}

// GetInvoices godoc
// @Summary List invoices for a store
// @Tags invoices
// @Param store_id query string true "Store ID"
// @Param status query string false "Filter by status"
// @Success 200 {array} models.Invoice
// @Router /invoices [get]
// @Security BearerAuth
func (h *InvoiceHandler) GetInvoices(c *gin.Context) {
	// Enforce store_id from JWT (data isolation)
	storeID, hasStore := c.Get("store_id")
	if !hasStore {
		// Super admin can query any store
		queryStoreID := c.Query("store_id")
		if queryStoreID != "" {
			storeID = queryStoreID
			hasStore = true
		}
	}

	if !hasStore {
		c.JSON(http.StatusBadRequest, gin.H{"error": "store_id is required"})
		return
	}

	invoices, err := h.invoiceRepo.FindByStoreID(storeID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch invoices"})
		return
	}

	c.JSON(http.StatusOK, invoices)
}

// CreateInvoice godoc
// @Summary Create a new invoice
// @Tags invoices
// @Param request body CreateInvoiceRequest true "Invoice data"
// @Success 201 {object} models.Invoice
// @Router /invoices [post]
// @Security BearerAuth
func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {
	var req CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	memberID, err := uuid.Parse(req.MemberID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member_id"})
		return
	}

	var dueDate *time.Time
	if req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", req.DueDate)
		if err == nil {
			dueDate = &parsed
		}
	}

	invoice := &models.Invoice{
		MemberID:    memberID,
		Amount:      req.Amount,
		Description: req.Description,
		Status:      models.InvoiceDraft,
		DueDate:     dueDate,
	}

	// Enforce store_id from JWT (data isolation), or allow super admin to pass in body
	storeID, hasStore := c.Get("store_id")
	if hasStore {
		parsedID, err := uuid.Parse(storeID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid store_id in token"})
			return
		}
		invoice.StoreID = parsedID
	} else if req.StoreID != "" {
		// Super admin: use store_id from request body
		parsedID, err := uuid.Parse(req.StoreID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid store_id"})
			return
		}
		invoice.StoreID = parsedID
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "store_id is required"})
		return
	}

	invoice.InvoiceNumber = h.invoiceRepo.GenerateInvoiceNumber(invoice.StoreID.String())

	if err := h.invoiceRepo.Create(invoice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create invoice"})
		return
	}

	c.JSON(http.StatusCreated, invoice)
}

// GetInvoice godoc
// @Summary Get invoice by ID
// @Tags invoices
// @Param id path string true "Invoice ID"
// @Success 200 {object} models.Invoice
// @Router /invoices/{id} [get]
// @Security BearerAuth
func (h *InvoiceHandler) GetInvoice(c *gin.Context) {
	id := c.Param("id")

	invoice, err := h.invoiceRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invoice not found"})
		return
	}

	c.JSON(http.StatusOK, invoice)
}

// UpdateInvoice godoc
// @Summary Update invoice
// @Tags invoices
// @Param id path string true "Invoice ID"
// @Param invoice body models.Invoice true "Updated data"
// @Success 200 {object} models.Invoice
// @Router /invoices/{id} [put]
// @Security BearerAuth
func (h *InvoiceHandler) UpdateInvoice(c *gin.Context) {
	id := c.Param("id")

	var updates models.Invoice
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invoice, err := h.invoiceRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invoice not found"})
		return
	}

	if updates.Amount > 0 {
		invoice.Amount = updates.Amount
	}
	if updates.Description != "" {
		invoice.Description = updates.Description
	}
	if updates.Status != "" {
		invoice.Status = updates.Status
	}

	if err := h.invoiceRepo.Update(invoice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update invoice"})
		return
	}

	c.JSON(http.StatusOK, invoice)
}

// DeleteInvoice godoc
// @Summary Delete invoice
// @Tags invoices
// @Param id path string true "Invoice ID"
// @Success 204
// @Router /invoices/{id} [delete]
// @Security BearerAuth
func (h *InvoiceHandler) DeleteInvoice(c *gin.Context) {
	id := c.Param("id")

	if err := h.invoiceRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete invoice"})
		return
	}

	c.Status(http.StatusNoContent)
}

// RecordPayment godoc
// @Summary Record payment for an invoice
// @Tags invoices
// @Param id path string true "Invoice ID"
// @Param request body RecordPaymentRequest true "Payment data"
// @Success 201 {object} models.Payment
// @Router /invoices/{id}/pay [post]
// @Security BearerAuth
func (h *InvoiceHandler) RecordPayment(c *gin.Context) {
	id := c.Param("id")

	var req RecordPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invoiceID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invoice id"})
		return
	}

	payment := &models.Payment{
		InvoiceID: invoiceID,
		Amount:    req.Amount,
		Method:    req.Method,
		Reference: req.Reference,
		Note:      req.Note,
		PaidAt:    time.Now(),
	}

	if payment.Method == "" {
		payment.Method = "manual"
	}

	if err := h.paymentRepo.Create(payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record payment"})
		return
	}

	// Update invoice status to paid
	invoice, err := h.invoiceRepo.FindByID(id)
	if err == nil {
		invoice.Status = models.InvoicePaid
		h.invoiceRepo.Update(invoice)
	}

	c.JSON(http.StatusCreated, payment)
}

// GetMemberInvoices godoc
// @Summary Get invoices for a specific member
// @Tags invoices
// @Param member_id query string true "Member ID"
// @Success 200 {array} models.Invoice
// @Router /invoices/member [get]
// @Security BearerAuth
func (h *InvoiceHandler) GetMemberInvoices(c *gin.Context) {
	memberID := c.Query("member_id")
	if memberID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "member_id is required"})
		return
	}

	invoices, err := h.invoiceRepo.FindByMemberID(memberID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch invoices"})
		return
	}

	c.JSON(http.StatusOK, invoices)
}
