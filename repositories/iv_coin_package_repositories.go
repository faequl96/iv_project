package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type IVCoinPackageRepositories interface {
	CreateIVCoinPackage(ivCoinPackage *models.IVCoinPackage) error
	GetIVCoinPackageByID(id uint) (*models.IVCoinPackage, error)
	GetIVCoinPackages() ([]models.IVCoinPackage, error)
	GetIVCoinPackagesByDiscountCategory(discountCategory string) ([]models.IVCoinPackage, error)
	UpdateIVCoinPackage(ivCoinPackage *models.IVCoinPackage) error
	DeleteIVCoinPackage(id uint) error
}

func IVCoinPackageRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateIVCoinPackage(ivCoinPackage *models.IVCoinPackage) error {
	return r.db.Create(ivCoinPackage).Error
}

func (r *repository) GetIVCoinPackageByID(id uint) (*models.IVCoinPackage, error) {
	var ivCoinPackage models.IVCoinPackage
	err := r.db.First(&ivCoinPackage, id).Error
	return &ivCoinPackage, err
}

func (r *repository) GetIVCoinPackages() ([]models.IVCoinPackage, error) {
	var ivCoinPackages []models.IVCoinPackage
	err := r.db.Find(&ivCoinPackages).Error
	return ivCoinPackages, err
}

func (r *repository) GetIVCoinPackagesByDiscountCategory(discountCategory string) ([]models.IVCoinPackage, error) {
	var ivCoinPackage []models.IVCoinPackage
	err := r.db.Preload("Category").Preload("Review").Find(&ivCoinPackage, "discount_category = ?", discountCategory).Error
	return ivCoinPackage, err
}

func (r *repository) UpdateIVCoinPackage(ivCoinPackage *models.IVCoinPackage) error {
	return r.db.Save(ivCoinPackage).Error
}

func (r *repository) DeleteIVCoinPackage(id uint) error {
	return r.db.Delete(&models.IVCoinPackage{}, id).Error
}
