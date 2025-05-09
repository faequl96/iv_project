package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type InvitationRepositories interface {
	CreateInvitation(invitation *models.Invitation) error
	GetInvitationByID(id uint) (*models.Invitation, error)
	GetInvitations() ([]models.Invitation, error)
	GetInvitationsByUserID(userID string) ([]models.Invitation, error)
	UpdateInvitation(invitation *models.Invitation) error
	DeleteInvitation(id uint) error
}

func InvitationRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateInvitation(invitation *models.Invitation) error {
	tx := r.db.Begin()
	if err := tx.Create(invitation).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *repository) GetInvitationByID(id uint) (*models.Invitation, error) {
	var invitation models.Invitation
	err := r.db.Preload("InvitationData.Gallery").Preload("InvitationTheme").First(&invitation, id).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (r *repository) GetInvitations() ([]models.Invitation, error) {
	var invitations []models.Invitation
	err := r.db.Preload("InvitationData.Gallery").Find(&invitations).Error
	return invitations, err
}

func (r *repository) GetInvitationsByUserID(userID string) ([]models.Invitation, error) {
	var invitations []models.Invitation
	err := r.db.Preload("InvitationData.Gallery").Find(&invitations, "user_id = ?", userID).Error
	return invitations, err
}

func (r *repository) UpdateInvitation(invitation *models.Invitation) error {
	tx := r.db.Begin()
	if err := tx.Save(invitation).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Save(invitation.InvitationData).Error; err != nil {
		tx.Rollback()
		return err
	}
	if invitation.InvitationData.Gallery != nil {
		if err := tx.Save(invitation.InvitationData.Gallery).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (r *repository) DeleteInvitation(id uint) error {
	tx := r.db.Begin()
	if err := tx.Delete(&models.Invitation{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
