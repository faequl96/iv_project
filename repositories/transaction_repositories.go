package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type TransactionRepositories interface {
	Create(transaction *models.Transaction) error
	FindByID(id uint) (*models.Transaction, error)
	FindByUserID(userID string) ([]models.Transaction, error)
	FindAll() ([]models.Transaction, error) // Tambahan
	Update(transaction *models.Transaction) error
}

func TransactionRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *repository) FindByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Invitation").Preload("IVCoinPackage").Where("id = ?", id).First(&transaction).Error
	return &transaction, err
}

func (r *repository) FindAll() ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("Invitation").Preload("IVCoinPackage").Find(&transactions).Error
	return transactions, err
}

func (r *repository) FindByUserID(userID string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("Invitation").Preload("IVCoinPackage").Where("user_id = ?", userID).Find(&transactions).Error
	return transactions, err
}

func (r *repository) Update(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}
