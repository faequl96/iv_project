package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type DiscountCategoryRepositories interface {
	CreateDiscountCategory(discountCategory *models.DiscountCategory) error
	GetDiscountCategoryByID(id uint) (*models.DiscountCategory, error)
	GetDiscountCategories() ([]models.DiscountCategory, error)
	GetDiscountCategoriesByIDs(ids []uint) ([]models.DiscountCategory, error)
	UpdateDiscountCategory(discountCategory *models.DiscountCategory) error
	DeleteDiscountCategory(id uint) error
}

func DiscountCategoryRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateDiscountCategory(discountCategory *models.DiscountCategory) error {
	return r.db.Create(discountCategory).Error
}

func (r *repository) GetDiscountCategoryByID(id uint) (*models.DiscountCategory, error) {
	var discountCategory models.DiscountCategory
	err := r.db.First(&discountCategory, id).Error
	if err != nil {
		return nil, err
	}
	return &discountCategory, nil
}

func (r *repository) GetDiscountCategories() ([]models.DiscountCategory, error) {
	var discountCategories []models.DiscountCategory
	err := r.db.Find(&discountCategories).Error
	return discountCategories, err
}

func (r *repository) GetDiscountCategoriesByIDs(ids []uint) ([]models.DiscountCategory, error) {
	var discountCategories []models.DiscountCategory
	err := r.db.Find(&discountCategories, "id IN ?", ids).Error
	return discountCategories, err
}

func (r *repository) UpdateDiscountCategory(discountCategory *models.DiscountCategory) error {
	return r.db.Save(discountCategory).Error
}

func (r *repository) DeleteDiscountCategory(id uint) error {
	return r.db.Delete(&models.DiscountCategory{}, id).Error
}
