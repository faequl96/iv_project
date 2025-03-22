package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type IVCoinPackageRepositories interface {
	CreateIVCoinPackage(ivCoinPackage *models.IVCoinPackage) error
	GetIVCoinPackageByID(id uint) (*models.IVCoinPackage, error)
	GetIVCoinPackages() ([]models.IVCoinPackage, error)
	GetIVCoinPackagesByDiscountCategoryID(discountCategoryID uint) ([]models.IVCoinPackage, error)
	UpdateIVCoinPackage(ivCoinPackage *models.IVCoinPackage) error
	DeleteIVCoinPackage(id uint) error
}

func IVCoinPackageRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateIVCoinPackage(ivCoinPackage *models.IVCoinPackage) error {
	tx := r.db.Begin()
	if err := tx.Create(ivCoinPackage).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *repository) GetIVCoinPackageByID(id uint) (*models.IVCoinPackage, error) {
	var ivCoinPackage models.IVCoinPackage
	err := r.db.Preload("DiscountCategories").First(&ivCoinPackage, id).Error
	return &ivCoinPackage, err
}

func (r *repository) GetIVCoinPackages() ([]models.IVCoinPackage, error) {
	var ivCoinPackages []models.IVCoinPackage
	err := r.db.Preload("DiscountCategories").Find(&ivCoinPackages).Error
	return ivCoinPackages, err
}

func (r *repository) GetIVCoinPackagesByDiscountCategoryID(discountCategoryID uint) ([]models.IVCoinPackage, error) {
	var ivCoinPackages []models.IVCoinPackage
	err := r.db.Preload("DiscountCategories").Where("id IN (?)", r.db.Table("iv_coin_package_discount_categories").
		Select("iv_coin_package_id").
		Where("discount_category_id = ?", discountCategoryID)).
		Find(&ivCoinPackages).Error
	if err != nil {
		return nil, err
	}
	return ivCoinPackages, nil
}

func (r *repository) UpdateIVCoinPackage(ivCoinPackage *models.IVCoinPackage) error {
	tx := r.db.Begin()
	if err := tx.Save(ivCoinPackage).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(ivCoinPackage).Association("DiscountCategories").Replace(ivCoinPackage.DiscountCategories); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *repository) DeleteIVCoinPackage(id uint) error {
	tx := r.db.Begin()
	if err := tx.Delete(&models.IVCoinPackage{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
