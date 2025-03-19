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
	return r.db.Create(invitationData).Error
}

func (r *repository) GetInvitationDataByID(id uint) (*models.InvitationData, error) {
	var invitationData models.InvitationData
	err := r.db.Preload("Gallery").First(&invitationData, id).Error
	return &invitationData, err
}

func (r *repository) UpdateInvitationData(invitationData *models.InvitationData) error {
	return r.db.Save(invitationData).Error
}

func (r *repository) DeleteInvitationData(id uint) error {
	return r.db.Delete(&models.InvitationData{}, id).Error
}
