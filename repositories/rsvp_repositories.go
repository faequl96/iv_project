package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type RSVPRepositories interface {
	CreateRSVP(rsvp *models.RSVP) error
	GetRSVPsByInvitationID(invitationID uint) ([]models.RSVP, error)
}

func RSVPRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateRSVP(rsvp *models.RSVP) error {
	return r.db.Create(rsvp).Error
}

func (r *repository) GetRSVPsByInvitationID(invitationID uint) ([]models.RSVP, error) {
	var rsvps []models.RSVP
	err := r.db.Find(&rsvps, "invitation_id = ?", invitationID).Error
	return rsvps, err
}
