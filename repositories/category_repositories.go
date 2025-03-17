package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type CategoryRepositories interface {
	CreateCategory(category *models.Category) error
	GetCategoryByID(id uint) (*models.Category, error)
	GetCategories() ([]models.Category, error)
	GetCategoriesByIDs(ids []uint) ([]models.Category, error)
	UpdateCategory(category *models.Category) error
	DeleteCategory(id uint) error
}

func CategoryRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateCategory(category *models.Category) error {
	tx := r.db.Begin()
	if err := tx.Create(category).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *repository) GetCategoryByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *repository) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *repository) GetCategoriesByIDs(ids []uint) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("id IN ?", ids).Find(&categories).Error
	return categories, err
}

func (r *repository) UpdateCategory(category *models.Category) error {
	tx := r.db.Begin()
	if err := tx.Save(category).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *repository) DeleteCategory(id uint) error {
	tx := r.db.Begin()
	if err := tx.Delete(&models.Category{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
