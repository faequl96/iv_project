package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type TransactionRepositories interface {
	CreateTransaction(transaction *models.Transaction) error
	GetTransactionByID(id uint) (*models.Transaction, error)
	GetTransactionsByUserID(userID string) ([]models.Transaction, error)
	GetTransactions() ([]models.Transaction, error)
	UpdateTransaction(transaction *models.Transaction) error
	DeleteTransaction(id uint) error
}

func TransactionRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateTransaction(transaction *models.Transaction) error {
	tx := r.db.Begin()
	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *repository) GetTransactionByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Invitation").Preload("IVCoinPackage").Where("id = ?", id).First(&transaction).Error
	return &transaction, err
}

func (r *repository) GetTransactions() ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("Invitation").Preload("IVCoinPackage").Find(&transactions).Error
	return transactions, err
}

func (r *repository) GetTransactionsByUserID(userID string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("Invitation").Preload("IVCoinPackage").Where("user_id = ?", userID).Find(&transactions).Error
	return transactions, err
}

func (r *repository) UpdateTransaction(transaction *models.Transaction) error {
	tx := r.db.Begin()
	if err := tx.Save(transaction).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *repository) DeleteTransaction(id uint) error {
	tx := r.db.Begin()
	if err := tx.Where("id = ?", id).Delete(&models.Transaction{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
