package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type InvitationThemeRepositories interface {
	CreateInvitationTheme(invitationTheme models.InvitationTheme) error
	GetInvitationThemeByID(invitationThemeID uint) (models.InvitationTheme, error)
	GetInvitationThemes() ([]models.InvitationTheme, error)
	GetInvitationThemesByCategory(category string) ([]models.InvitationTheme, error)
	UpdateInvitationTheme(invitationTheme models.InvitationTheme) error
	DeleteInvitationTheme(invitationTheme models.InvitationTheme) error
}

func InvitationThemeRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateInvitationTheme(invitationTheme models.InvitationTheme) error {
	return r.db.Create(&invitationTheme).Error
}

func (r *repository) GetInvitationThemeByID(invitationThemeID uint) (models.InvitationTheme, error) {
	var invitationTheme models.InvitationTheme
	err := r.db.Where("id = ?", invitationThemeID).First(&invitationTheme).Error
	return invitationTheme, err
}

func (r *repository) GetInvitationThemes() ([]models.InvitationTheme, error) {
	var invitationThemes []models.InvitationTheme
	err := r.db.Find(&invitationThemes).Error
	return invitationThemes, err
}

func (r *repository) GetInvitationThemesByCategory(category string) ([]models.InvitationTheme, error) {
	var invitationThemes []models.InvitationTheme
	err := r.db.Where("category = ?", category).Find(&invitationThemes).Error
	return invitationThemes, err
}

func (r *repository) UpdateInvitationTheme(invitationTheme models.InvitationTheme) error {
	return r.db.Save(&invitationTheme).Error
}

func (r *repository) DeleteInvitationTheme(invitationTheme models.InvitationTheme) error {
	return r.db.Delete(&invitationTheme).Error
}
