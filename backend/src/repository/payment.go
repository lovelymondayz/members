package repository

import (
	"github.com/lovelymondayz/members/backend/src/config"
	"github.com/lovelymondayz/members/backend/src/models"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{db: config.DB}
}

func (r *PaymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *PaymentRepository) FindByInvoiceID(invoiceID string) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Where("invoice_id = ?", invoiceID).Order("paid_at DESC").Find(&payments).Error
	return payments, err
}

func (r *PaymentRepository) FindByStoreID(storeID string) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Joins("JOIN invoices ON payments.invoice_id = invoices.invoice_id").
		Where("invoices.store_id = ?", storeID).
		Order("payments.paid_at DESC").
		Find(&payments).Error
	return payments, err
}
