package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type InvitationThemeRepositories interface {
	CreateInvitationTheme(invitationTheme *models.InvitationTheme) error
	GetInvitationThemeByID(id uint) (*models.InvitationTheme, error)
	GetInvitationThemes() ([]models.InvitationTheme, error)
	GetInvitationThemesByCategory(category string) ([]models.InvitationTheme, error)
	UpdateInvitationTheme(invitationTheme *models.InvitationTheme) error
	DeleteInvitationTheme(id uint) error
}

func InvitationThemeRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateInvitationTheme(invitationTheme *models.InvitationTheme) error {
	tx := r.db.Begin()
	if err := tx.Create(invitationTheme).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *repository) GetInvitationThemeByID(id uint) (*models.InvitationTheme, error) {
	var invitationTheme models.InvitationTheme
	err := r.db.Preload("Category").Preload("Review").First(&invitationTheme, id).Error
	return &invitationTheme, err
}

func (r *repository) GetInvitationThemes() ([]models.InvitationTheme, error) {
	var invitationThemes []models.InvitationTheme
	err := r.db.Preload("Category").Preload("Review").Find(&invitationThemes).Error
	return invitationThemes, err
}

func (r *repository) GetInvitationThemesByCategory(category string) ([]models.InvitationTheme, error) {
	var invitationThemes []models.InvitationTheme
	err := r.db.Preload("Category").Preload("Review").Find(&invitationThemes, "category = ?", category).Error
	return invitationThemes, err
}

func (r *repository) UpdateInvitationTheme(invitationTheme *models.InvitationTheme) error {
	tx := r.db.Begin()
	if err := tx.Save(invitationTheme).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *repository) DeleteInvitationTheme(id uint) error {
	tx := r.db.Begin()
	if err := tx.Delete(&models.InvitationTheme{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
