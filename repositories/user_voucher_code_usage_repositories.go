package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type UserVoucherCodeUsageRepositories interface {
	CreateUserVoucherCodeUsage(userVoucherCodeUsage *models.UserVoucherCodeUsage) error
	GetUserVoucherCodeUsageByUserAndVoucherCodeID(userId string, voucherCodeId uint) (*models.UserVoucherCodeUsage, error)
	UpdateUserVoucherCodeUsage(userVoucherCodeUsage *models.UserVoucherCodeUsage) error
}

func UserVoucherCodeUsageRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateUserVoucherCodeUsage(userVoucherCodeUsage *models.UserVoucherCodeUsage) error {
	return r.db.Create(userVoucherCodeUsage).Error
}

func (r *repository) GetUserVoucherCodeUsageByUserAndVoucherCodeID(userId string, voucherCodeId uint) (*models.UserVoucherCodeUsage, error) {
	var userVoucherCodeUsage models.UserVoucherCodeUsage
	err := r.db.Where("user_id = ? AND voucher_code_id = ?", userId, voucherCodeId).First(&userVoucherCodeUsage).Error
	if err != nil {
		return nil, err
	}
	return &userVoucherCodeUsage, nil
}

func (r *repository) UpdateUserVoucherCodeUsage(userVoucherCodeUsage *models.UserVoucherCodeUsage) error {
	return r.db.Save(userVoucherCodeUsage).Error
}
