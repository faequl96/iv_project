package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type InvitationDataRepositories interface {
	CreateInvitationData(invitationData *models.InvitationData) error
	GetInvitationDataByID(id uint) (*models.InvitationData, error)
	UpdateInvitationData(invitationData *models.InvitationData) error
	DeleteInvitationData(id uint) error
}

func InvitationDataRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateInvitationData(invitationData *models.InvitationData) error {
	tx := r.db.Begin()
	if err := tx.Create(invitationData).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *repository) GetInvitationDataByID(id uint) (*models.InvitationData, error) {
	var invitationData models.InvitationData
	err := r.db.Preload("Gallery").Where("id = ?", id).First(&invitationData).Error
	return &invitationData, err
}

func (r *repository) UpdateInvitationData(invitationData *models.InvitationData) error {
	tx := r.db.Begin()
	if err := tx.Save(invitationData).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *repository) DeleteInvitationData(id uint) error {
	tx := r.db.Begin()
	if err := tx.Where("id = ?", id).Delete(&models.InvitationData{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
