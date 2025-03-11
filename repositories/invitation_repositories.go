package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type InvitationRepositories interface {
	CreateInvitation(invitation models.Invitation) error
	GetInvitationByID(invitationID uint) (models.Invitation, error)
	GetInvitationsByUserID(userID uint) ([]models.Invitation, error)
}

func InvitationRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateInvitation(invitation models.Invitation) error {
	return r.db.Create(invitation).Error
}

func (r *repository) GetInvitationByID(invitationID uint) (models.Invitation, error) {
	var invitation models.Invitation
	err := r.db.Preload("InvitationData").Where("id = ?", invitationID).First(&invitation).Error
	return invitation, err
}

func (r *repository) GetInvitationsByUserID(userID uint) ([]models.Invitation, error) {
	var invitations []models.Invitation
	err := r.db.Preload("InvitationData").Where("user_id = ?", userID).Find(&invitations).Error
	return invitations, err
}
