package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type InvitationDataRepositories interface {
	CreateInvitationData(invitationData models.InvitationData) error
	GetInvitationDataByID(invitationDataID uint) (models.InvitationData, error)
	UpdateInvitationData(invitationData models.InvitationData) error
	DeleteInvitationData(invitationData models.InvitationData) error
}

func InvitationDataRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateInvitationData(invitationData models.InvitationData) error {
	tx := r.db.Begin()

	// Update the Invitation in the database
	if err := tx.Create(&invitationData).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction if everything is successful
	return tx.Commit().Error
}

func (r *repository) GetInvitationDataByID(invitationDataID uint) (models.InvitationData, error) {
	var invitationData models.InvitationData
	err := r.db.Preload("Gallery").Where("id = ?", invitationDataID).First(&invitationData).Error
	return invitationData, err
}

func (r *repository) UpdateInvitationData(invitationData models.InvitationData) error {
	tx := r.db.Begin()

	// Update the Invitation in the database
	if err := tx.Create(&invitationData).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction if everything is successful
	return tx.Commit().Error
}

func (r *repository) DeleteInvitationData(invitationData models.InvitationData) error {
	return r.db.Delete(&invitationData).Error
}
