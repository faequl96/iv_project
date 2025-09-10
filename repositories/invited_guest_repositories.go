package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type InvitedGuestRepositories interface {
	CreateInvitedGuest(invitedGuest *models.InvitedGuest) error
	GetInvitedGuestByID(id uint) (*models.InvitedGuest, error)
	GetInvitedGuestsByInvitationID(invitationID uint) ([]models.InvitedGuest, error)
	UpdateInvitedGuest(invitedGuest *models.InvitedGuest) error
}

func InvitedGuestRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateInvitedGuest(invitedGuest *models.InvitedGuest) error {
	return r.db.Create(invitedGuest).Error
}

func (r *repository) GetInvitedGuestByID(id uint) (*models.InvitedGuest, error) {
	var invitedGuest models.InvitedGuest
	err := r.db.First(&invitedGuest, id).Error
	if err != nil {
		return nil, err
	}
	return &invitedGuest, nil
}

func (r *repository) GetInvitedGuestsByInvitationID(invitationID uint) ([]models.InvitedGuest, error) {
	var invitedGuests []models.InvitedGuest
	err := r.db.Find(&invitedGuests, "invitation_id = ?", invitationID).Error
	return invitedGuests, err
}

func (r *repository) UpdateInvitedGuest(invitedGuest *models.InvitedGuest) error {
	return r.db.Save(invitedGuest).Error
}
