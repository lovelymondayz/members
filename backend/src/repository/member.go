package repository

import (
	"github.com/lovelymondayz/members/backend/src/config"
	"github.com/lovelymondayz/members/backend/src/models"
	"gorm.io/gorm"
)

type MemberRepository struct {
	db *gorm.DB
}

func NewMemberRepository() *MemberRepository {
	return &MemberRepository{db: config.DB}
}

func (r *MemberRepository) FindByStoreID(storeID string) ([]models.Member, error) {
	var members []models.Member
	err := r.db.Where("store_id = ?", storeID).Find(&members).Error
	// Load related data separately to avoid circular FK issues
	for i := range members {
		if members[i].UserID != nil {
			var user models.User
			r.db.First(&user, "user_id = ?", *members[i].UserID)
			members[i].UserName = user.Name
		}
		var store models.Store
		r.db.First(&store, "store_id = ?", members[i].StoreID)
		members[i].StoreName = store.Name
		members[i].StoreCardColor = store.CardColorHex
	}
	return members, err
}

func (r *MemberRepository) FindAll() ([]models.Member, error) {
	var members []models.Member
	err := r.db.Order("created_at desc").Find(&members).Error
	for i := range members {
		if members[i].UserID != nil {
			var user models.User
			r.db.First(&user, "user_id = ?", *members[i].UserID)
			members[i].UserName = user.Name
		}
		var store models.Store
		r.db.First(&store, "store_id = ?", members[i].StoreID)
		members[i].StoreName = store.Name
		members[i].StoreCardColor = store.CardColorHex
	}
	return members, err
}

func (r *MemberRepository) FindByID(id string) (*models.Member, error) {
	var member models.Member
	err := r.db.First(&member, "member_id = ?", id).Error
	return &member, err
}

func (r *MemberRepository) FindByIDWithRelations(id string) (*models.Member, error) {
	var member models.Member
	err := r.db.First(&member, "member_id = ?", id).Error
	if err != nil {
		return &member, err
	}
	// Load transient relation data
	if member.UserID != nil {
		var user models.User
		r.db.First(&user, "user_id = ?", *member.UserID)
		member.UserName = user.Name
	}
	var store models.Store
	r.db.First(&store, "store_id = ?", member.StoreID)
	member.StoreName = store.Name
	member.StoreCardColor = store.CardColorHex
	return &member, nil
}

func (r *MemberRepository) Create(member *models.Member) error {
	return r.db.Create(member).Error
}

func (r *MemberRepository) Update(member *models.Member) error {
	return r.db.Save(member).Error
}

func (r *MemberRepository) Delete(id string) error {
	return r.db.Delete(&models.Member{}, "member_id = ?", id).Error
}
