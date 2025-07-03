package repositories

import (
	query_dto "iv_project/dto/query"
	"iv_project/models"
	"iv_project/pkg/utils"

	"gorm.io/gorm"
)

type TransactionRepositories interface {
	CreateTransaction(transaction *models.Transaction) error
	GetTransactionByID(id string) (*models.Transaction, error)
	GetTransactionByReferenceNumber(referenceNumber string) (*models.Transaction, error)
	GetTransactions(query *query_dto.QueryRequest) ([]models.Transaction, error)
	GetTransactionsByUserID(userID string) ([]models.Transaction, error)
	UpdateTransaction(transaction *models.Transaction) error
	DeleteTransaction(id string) error
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

func (r *repository) GetTransactionByID(id string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Where("id = ?", id).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *repository) GetTransactionByReferenceNumber(referenceNumber string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Where("reference_number = ?", referenceNumber).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *repository) GetTransactions(query *query_dto.QueryRequest) ([]models.Transaction, error) {
	var transactions []models.Transaction
	db := utils.ImplementQuery(r.db, query)
	err := db.Find(&transactions).Error
	return transactions, err
}

func (r *repository) GetTransactionsByUserID(userID string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Order("created_at DESC").Find(&transactions, "user_id = ?", userID).Error
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

func (r *repository) DeleteTransaction(id string) error {
	tx := r.db.Begin()
	if err := tx.Where("id = ?", id).Delete(&models.Transaction{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
