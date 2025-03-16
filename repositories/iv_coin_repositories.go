package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type IVCoinRepositories interface {
	GetIVCoinByID(id uint) (*models.IVCoin, error)
	UpdateIVCoin(ivCoin *models.IVCoin) error
}

func IVCoinRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetIVCoinByID(id uint) (*models.IVCoin, error) {
	var ivCoin *models.IVCoin
	err := r.db.Where("id = ?", id).First(&ivCoin).Error
	return ivCoin, err
}

func (r *repository) UpdateIVCoin(ivCoin *models.IVCoin) error {
	return r.db.Save(&ivCoin).Error
}
