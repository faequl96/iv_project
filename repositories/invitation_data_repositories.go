package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type InvitationDataRepositories interface {
	GetInvitationDataByID(invitationDataID uint) (models.InvitationData, error)
	UpdateInvitationData(invitationData models.InvitationData) error
}

func InvitationDataRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetInvitationDataByID(invitationDataID uint) (models.InvitationData, error) {
	var invitationData models.InvitationData
	err := r.db.Where("id = ?", invitationDataID).First(&invitationData).Error
	return invitationData, err
}

func (r *repository) UpdateInvitationData(invitationData models.InvitationData) error {
	err := r.db.Save(&invitationData).Error
	return err
}
