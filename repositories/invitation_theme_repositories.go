package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type InvitationThemeRepositories interface {
	CreateInvitationTheme(invitationTheme *models.InvitationTheme) error
	GetInvitationThemeByID(id uint) (*models.InvitationTheme, error)
	GetInvitationThemes() ([]models.InvitationTheme, error)
	GetInvitationThemesByCategoryID(categoryID uint) ([]models.InvitationTheme, error)
	GetInvitationThemesByDiscountCategoryID(discountCategoryID uint) ([]models.InvitationTheme, error)
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
	err := r.db.Preload("Categories").Preload("Reviews").First(&invitationTheme, id).Error
	return &invitationTheme, err
}

func (r *repository) GetInvitationThemes() ([]models.InvitationTheme, error) {
	var invitationThemes []models.InvitationTheme
	err := r.db.Preload("Categories").Preload("Reviews").Find(&invitationThemes).Error
	return invitationThemes, err
}

func (r *repository) GetInvitationThemesByCategoryID(categoryID uint) ([]models.InvitationTheme, error) {
	var invitationThemes []models.InvitationTheme
	err := r.db.Preload("Categories").Preload("Reviews").Where("id IN (?)", r.db.Table("invitation_theme_categories").
		Select("invitation_theme_id").
		Where("category_id = ?", categoryID)).
		Find(&invitationThemes).Error
	if err != nil {
		return nil, err
	}
	return invitationThemes, nil
}

func (r *repository) GetInvitationThemesByDiscountCategoryID(discountCategoryID uint) ([]models.InvitationTheme, error) {
	var invitationThemes []models.InvitationTheme
	err := r.db.Preload("Categories").Preload("Reviews").Where("id IN (?)", r.db.Table("invitation_theme_discount_categories").
		Select("invitation_theme_id").
		Where("discount_category_id = ?", discountCategoryID)).
		Find(&invitationThemes).Error
	if err != nil {
		return nil, err
	}
	return invitationThemes, nil
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
