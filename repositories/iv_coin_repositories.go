package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

// IVCoinRepositories defines the interface for performing CRUD operations on the IVCoin model in the database.
type IVCoinRepositories interface {
	GetIVCoinByID(id uint) (*models.IVCoin, error) // Fetch IVCoin by its ID.
	UpdateIVCoin(ivCoin *models.IVCoin) error      // Update IVCoin information.
}

// IVCoinRepository initializes the repository with a database connection.
func IVCoinRepository(db *gorm.DB) *repository {
	return &repository{db} // Return a repository instance initialized with the database.
}

// GetIVCoinByID retrieves an IVCoin by its ID.
func (r *repository) GetIVCoinByID(id uint) (*models.IVCoin, error) {
	var ivCoin models.IVCoin // Use value instead of pointer to avoid potential nil pointer issues.
	err := r.db.Where("id = ?", id).First(&ivCoin).Error
	if err != nil {
		return nil, err // Return nil and the error if the IVCoin is not found.
	}
	return &ivCoin, nil // Return the found IVCoin.
}

// UpdateIVCoin updates the information of an existing IVCoin in the database.
// `Save` performs an `INSERT` if the IVCoin doesn't exist, or an `UPDATE` if it does.
// Using a pointer (`&ivCoin`) to modify the original instance directly.
func (r *repository) UpdateIVCoin(ivCoin *models.IVCoin) error {
	// Use `Save` for both insert and update operations, as it checks if the record exists.
	return r.db.Save(ivCoin).Error // Return any error that occurs during the operation.
}
