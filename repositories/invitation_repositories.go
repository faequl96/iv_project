package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type InvitationRepositories interface {
	CreateInvitation(invitation models.Invitation) error
	GetInvitationByID(invitationID uint) (models.Invitation, error)
	GetInvitations() ([]models.Invitation, error)
	GetInvitationsByUserID(userID string) ([]models.Invitation, error)
	UpdateInvitation(invitation models.Invitation) error
	DeleteInvitation(invitation models.Invitation) error
}

func InvitationRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateInvitation(invitation models.Invitation) error {
	tx := r.db.Begin()

	// Update the Invitation in the database
	if err := tx.Create(&invitation).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction if everything is successful
	return tx.Commit().Error
}

func (r *repository) GetInvitationByID(invitationID uint) (models.Invitation, error) {
	var invitation models.Invitation
	err := r.db.Where("id = ?", invitationID).First(&invitation).Error
	return invitation, err
}

func (r *repository) GetInvitations() ([]models.Invitation, error) {
	var invitations []models.Invitation
	err := r.db.Find(&invitations).Error
	return invitations, err
}

func (r *repository) GetInvitationsByUserID(userID string) ([]models.Invitation, error) {
	var invitations []models.Invitation
	err := r.db.Preload("InvitationData").Where("user_id = ?", userID).Find(&invitations).Error
	return invitations, err
}

func (r *repository) UpdateInvitation(invitation models.Invitation) error {
	tx := r.db.Begin()

	// Update the Invitation in the database
	if err := tx.Save(&invitation).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction if everything is successful
	return tx.Commit().Error
}

func (r *repository) DeleteInvitation(invitation models.Invitation) error {
	return r.db.Delete(&invitation).Error
}
