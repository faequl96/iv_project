package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

// UserRepositories defines the interface for user-related database operations.
type UserRepositories interface {
	CreateUser(user *models.User) error
	GetUserByID(id string) (*models.User, error)
	GetUsers() ([]models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id string) error
}

// UserRepository initializes the repository with a database connection.
func UserRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// CreateUser inserts a new user into the database using a transaction.
// If an error occurs, the transaction is rolled back to maintain data integrity.
func (r *repository) CreateUser(user *models.User) error {
	tx := r.db.Begin()
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback() // Revert any changes if an error occurs.
		return err
	}
	return tx.Commit().Error // Finalize the transaction.
}

// GetUserByID retrieves a user by their ID with the IVCoin relationship preloaded.
func (r *repository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("IVCoin").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err // Return nil if no user is found or an error occurs.
	}
	return &user, nil
}

// GetUsers retrieves all users with the IVCoin relationship preloaded.
func (r *repository) GetUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("IVCoin").Find(&users).Error
	return users, err
}

// UpdateUser updates an existing user's information using a transaction.
// If an error occurs, the transaction is rolled back to prevent partial updates.
func (r *repository) UpdateUser(user *models.User) error {
	tx := r.db.Begin()
	if err := tx.Save(user).Error; err != nil {
		tx.Rollback() // Revert changes if something goes wrong.
		return err
	}
	return tx.Commit().Error // Finalize the transaction.
}

// DeleteUser removes a user by ID using a transaction to ensure consistency.
// If an error occurs, the transaction is rolled back to prevent data inconsistency.
func (r *repository) DeleteUser(id string) error {
	tx := r.db.Begin()
	if err := tx.Where("id = ?", id).Delete(&models.User{}).Error; err != nil {
		tx.Rollback() // Ensure no unintended deletions happen.
		return err
	}
	return tx.Commit().Error // Finalize the deletion.
}
