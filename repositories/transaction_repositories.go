package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type TransactionRepositories interface {
	CreateTransaction(transaction *models.Transaction) error
	GetTransactionByID(id uint) (*models.Transaction, error)
	GetTransactions() ([]models.Transaction, error)
	GetTransactionsByUserID(userID string) ([]models.Transaction, error)
	GetTransactionByReferenceNumber(referenceNumber string) (*models.Transaction, error)
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
	err := r.db.First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *repository) GetTransactions() ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Find(&transactions).Error
	return transactions, err
}

func (r *repository) GetTransactionsByUserID(userID string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Find(&transactions, "user_id = ?", userID).Error
	return transactions, err
}

func (r *repository) GetTransactionByReferenceNumber(referenceNumber string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Where("reference_number = ?", referenceNumber).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
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
	if err := tx.Delete(&models.Transaction{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
