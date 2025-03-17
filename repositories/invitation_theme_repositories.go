package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

// InvitationThemeRepositories defines the interface for performing CRUD operations
type InvitationThemeRepositories interface {
	CreateInvitationTheme(invitationTheme *models.InvitationTheme) error
	GetInvitationThemeByID(id uint) (*models.InvitationTheme, error)
	GetInvitationThemes() ([]models.InvitationTheme, error)
	GetInvitationThemesByCategory(category string) ([]models.InvitationTheme, error)
	UpdateInvitationTheme(invitationTheme *models.InvitationTheme) error
	DeleteInvitationTheme(id uint) error
}

// InvitationThemeRepository initializes the repository with a database connection.
func InvitationThemeRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// CreateInvitationTheme creates a new InvitationTheme in the database.
// Uses a transaction to ensure atomicity. If an error occurs, the transaction is rolled back.
func (r *repository) CreateInvitationTheme(invitationTheme *models.InvitationTheme) error {
	tx := r.db.Begin() // Start a new transaction.
	if err := tx.Create(invitationTheme).Error; err != nil {
		tx.Rollback() // Rollback the transaction if an error occurs.
		return err
	}
	return tx.Commit().Error // Commit the transaction if no errors occur.
}

// GetInvitationThemeByID retrieves an InvitationTheme by its ID, preloading the associated 'Review' relationship.
func (r *repository) GetInvitationThemeByID(id uint) (*models.InvitationTheme, error) {
	var invitationTheme models.InvitationTheme
	err := r.db.Preload("Review").First(&invitationTheme, id).Error
	if err != nil {
		return nil, err // Return nil and the error if no InvitationTheme is found.
	}
	return &invitationTheme, nil
}

// GetInvitationThemes retrieves all InvitationThemes, preloading the associated 'Review' relationship.
func (r *repository) GetInvitationThemes() ([]models.InvitationTheme, error) {
	var invitationThemes []models.InvitationTheme
	err := r.db.Preload("Review").Find(&invitationThemes).Error
	return invitationThemes, err // Return the list of InvitationThemes or an error if it occurs.
}

// GetInvitationThemesByCategory retrieves all InvitationThemes for a specific category.
func (r *repository) GetInvitationThemesByCategory(category string) ([]models.InvitationTheme, error) {
	var invitationThemes []models.InvitationTheme
	err := r.db.Preload("Review").Where("category = ?", category).Find(&invitationThemes).Error
	return invitationThemes, err // Return the list of InvitationThemes or an error if it occurs.
}

// UpdateInvitationTheme updates an existing InvitationTheme in the database.
// Uses a transaction to ensure atomicity. If an error occurs, the transaction is rolled back.
func (r *repository) UpdateInvitationTheme(invitationTheme *models.InvitationTheme) error {
	tx := r.db.Begin() // Start a new transaction.
	if err := tx.Save(invitationTheme).Error; err != nil {
		tx.Rollback() // Rollback the transaction if an error occurs.
		return err
	}
	return tx.Commit().Error // Commit the transaction if no errors occur.
}

// DeleteInvitationTheme deletes an InvitationTheme by its ID from the database.
// Uses a transaction to ensure atomicity. If an error occurs, the transaction is rolled back.
func (r *repository) DeleteInvitationTheme(id uint) error {
	tx := r.db.Begin() // Start a new transaction.
	if err := tx.Delete(&models.InvitationTheme{}, id).Error; err != nil {
		tx.Rollback() // Rollback the transaction if an error occurs.
		return err
	}
	return tx.Commit().Error // Commit the transaction if no errors occur.
}
