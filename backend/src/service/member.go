package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/lovelymondayz/members/backend/src/models"
	"github.com/lovelymondayz/members/backend/src/repository"
)

type MemberService struct {
	memberRepo *repository.MemberRepository
	userRepo   *repository.UserRepository
}

func NewMemberService(memberRepo *repository.MemberRepository, userRepo *repository.UserRepository) *MemberService {
	return &MemberService{
		memberRepo: memberRepo,
		userRepo:   userRepo,
	}
}

func (s *MemberService) GetMembersByStoreID(storeID string) ([]models.Member, error) {
	return s.memberRepo.FindByStoreID(storeID)
}

func (s *MemberService) GetMemberByID(id string) (*models.Member, error) {
	return s.memberRepo.FindByIDWithRelations(id)
}

func (s *MemberService) CreateMember(storeID, name, email, memberCode, tier string) (*models.Member, error) {
	var userID *uuid.UUID
	if email != "" {
		existingUser, err := s.userRepo.FindByEmail(email)
		if err == nil {
			userID = &existingUser.UserID
		}
	}

	parsedStoreID, err := uuid.Parse(storeID)
	if err != nil {
		return nil, fmt.Errorf("invalid store_id: %w", err)
	}

	member := &models.Member{
		StoreID:    parsedStoreID,
		MemberCode: memberCode,
		Tier:       tier,
	}

	if userID != nil {
		member.UserID = userID
	}

	if err := s.memberRepo.Create(member); err != nil {
		return nil, fmt.Errorf("failed to create member: %w", err)
	}

	return member, nil
}
