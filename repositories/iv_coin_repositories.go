package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type IVCoinRepositories interface {
	GetIVCoinByID(id uint) (*models.IVCoin, error)
	GetIVCoinByUserID(userID string) (*models.IVCoin, error)
	UpdateIVCoin(ivCoin *models.IVCoin) error
}

func IVCoinRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetIVCoinByID(id uint) (*models.IVCoin, error) {
	var ivCoin models.IVCoin
	err := r.db.First(&ivCoin, id).Error
	if err != nil {
		return nil, err
	}
	return &ivCoin, nil
}

func (r *repository) GetIVCoinByUserID(userID string) (*models.IVCoin, error) {
	var ivCoin models.IVCoin
	err := r.db.Where("user_id = ?", userID).First(&ivCoin).Error
	if err != nil {
		return nil, err
	}
	return &ivCoin, nil
}

func (r *repository) UpdateIVCoin(ivCoin *models.IVCoin) error {
	return r.db.Save(ivCoin).Error
}
