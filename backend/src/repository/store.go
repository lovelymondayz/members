package repository

import (
	"github.com/lovelymondayz/members/backend/src/config"
	"github.com/lovelymondayz/members/backend/src/models"
	"gorm.io/gorm"
)

type StoreRepository struct {
	db *gorm.DB
}

func NewStoreRepository() *StoreRepository {
	return &StoreRepository{db: config.DB}
}

func (r *StoreRepository) FindByAdminID(adminID string) (*models.Store, error) {
	var store models.Store
	err := r.db.Where("admin_id = ?", adminID).First(&store).Error
	return &store, err
}

func (r *StoreRepository) FindByID(id string) (*models.Store, error) {
	var store models.Store
	err := r.db.First(&store, "store_id = ?", id).Error
	return &store, err
}

func (r *StoreRepository) Create(store *models.Store) error {
	return r.db.Create(store).Error
}

func (r *StoreRepository) Update(store *models.Store) error {
	return r.db.Save(store).Error
}
