package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type VoucherCodeRepositories interface {
	CreateVoucherCode(voucherCode *models.VoucherCode) error
	GetVoucherCodeByID(id uint) (*models.VoucherCode, error)
	UpdateVoucherCode(voucherCode *models.VoucherCode) error
	DeleteVoucherCode(id uint) error
}

func VoucherCodeRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateVoucherCode(voucherCode *models.VoucherCode) error {
	return r.db.Create(voucherCode).Error
}

func (r *repository) GetVoucherCodeByID(id uint) (*models.VoucherCode, error) {
	var voucherCode models.VoucherCode
	err := r.db.First(&voucherCode, id).Error
	if err != nil {
		return nil, err
	}
	return &voucherCode, nil
}

func (r *repository) UpdateVoucherCode(voucherCode *models.VoucherCode) error {
	return r.db.Save(voucherCode).Error
}

func (r *repository) DeleteVoucherCode(id uint) error {
	return r.db.Delete(&models.VoucherCode{}, id).Error
}
