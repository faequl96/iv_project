package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type InvitationDataRepositories interface {
	GetInvitationDataByID(id uint) (*models.InvitationData, error)
	UpdateInvitationData(invitationData *models.InvitationData) error
}

func InvitationDataRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetInvitationDataByID(id uint) (*models.InvitationData, error) {
	var invitationData models.InvitationData
	err := r.db.Preload("Gallery").First(&invitationData, id).Error
	if err != nil {
		return nil, err
	}
	return &invitationData, nil
}

func (r *repository) UpdateInvitationData(invitationData *models.InvitationData) error {
	tx := r.db.Begin()
	if err := tx.Save(invitationData).Error; err != nil {
		tx.Rollback()
		return err
	}
	if invitationData.Gallery != nil {
		if err := tx.Save(invitationData.Gallery).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}
