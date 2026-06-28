package repository

import (
	"github.com/lovelymondayz/members/backend/src/config"
	"github.com/lovelymondayz/members/backend/src/models"
	"gorm.io/gorm"
)

type InvoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository() *InvoiceRepository {
	return &InvoiceRepository{db: config.DB}
}

// Load relations with user_name for find by store
func (r *InvoiceRepository) FindByStoreID(storeID string) ([]models.Invoice, error) {
	var invoices []models.Invoice
	err := r.db.Where("store_id = ?", storeID).Order("created_at DESC").Find(&invoices).Error
	if err != nil {
		return invoices, err
	}
	// Populate member info (user_name) for each invoice
	for i := range invoices {
		var member models.Member
		r.db.First(&member, "member_id = ?", invoices[i].MemberID)
		if member.UserID != nil {
			var user models.User
			r.db.First(&user, "user_id = ?", *member.UserID)
			invoices[i].MemberName = user.Name
		}
		invoices[i].MemberCode = member.MemberCode
	}
	return invoices, nil
}

func (r *InvoiceRepository) FindByID(id string) (*models.Invoice, error) {
	var invoice models.Invoice
	err := r.db.First(&invoice, "invoice_id = ?", id).Error
	if err != nil {
		return &invoice, err
	}
	// Load member info
	var member models.Member
	r.db.First(&member, "member_id = ?", invoice.MemberID)
	if member.UserID != nil {
		var user models.User
		r.db.First(&user, "user_id = ?", *member.UserID)
		invoice.MemberName = user.Name
	}
	invoice.MemberCode = member.MemberCode
	// Load payments
	var payments []models.Payment
	r.db.Where("invoice_id = ?", id).Order("paid_at DESC").Find(&payments)
	invoice.Payments = payments
	return &invoice, nil
}

func (r *InvoiceRepository) FindByMemberID(memberID string) ([]models.Invoice, error) {
	var invoices []models.Invoice
	err := r.db.Where("member_id = ?", memberID).Order("created_at DESC").Find(&invoices).Error
	return invoices, err
}

func (r *InvoiceRepository) FindAll() ([]models.Invoice, error) {
	var invoices []models.Invoice
	err := r.db.Order("created_at DESC").Find(&invoices).Error
	if err != nil {
		return invoices, err
	}
	for i := range invoices {
		var member models.Member
		r.db.First(&member, "member_id = ?", invoices[i].MemberID)
		if member.UserID != nil {
			var user models.User
			r.db.First(&user, "user_id = ?", *member.UserID)
			invoices[i].MemberName = user.Name
		}
		invoices[i].MemberCode = member.MemberCode
	}
	return invoices, nil
}

func (r *InvoiceRepository) Create(invoice *models.Invoice) error {
	return r.db.Create(invoice).Error
}

func (r *InvoiceRepository) Update(invoice *models.Invoice) error {
	return r.db.Save(invoice).Error
}

func (r *InvoiceRepository) Delete(id string) error {
	return r.db.Delete(&models.Invoice{}, "invoice_id = ?", id).Error
}

func (r *InvoiceRepository) GenerateInvoiceNumber(storeID string) string {
	var count int64
	r.db.Model(&models.Invoice{}).Where("store_id = ?", storeID).Count(&count)
	return generateInvoiceNumber(storeID, count+1)
}

func generateInvoiceNumber(storeID string, seq int64) string {
	// Simple format: INV-{last4 of storeID}-{sequence}
	if len(storeID) > 4 {
		storeID = storeID[len(storeID)-4:]
	}
	return "INV-" + storeID + "-" + padInt(seq, 4)
}

func padInt(n int64, width int) string {
	s := ""
	for i := 0; i < width; i++ {
		s = "0" + s
	}
	result := s + string(rune('0'+n%10))
	return result[len(result)-width:]
}
