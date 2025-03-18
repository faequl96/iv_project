package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type GalleryRepositories interface {
	GetGalleryByID(id uint) (*models.Gallery, error)
	DeleteGallery(id uint) error
}

func GalleryRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetGalleryByID(id uint) (*models.Gallery, error) {
	var Gallery models.Gallery
	err := r.db.Where("id = ?", id).First(&Gallery).Error
	return &Gallery, err
}

func (r *repository) DeleteGallery(id uint) error {
	tx := r.db.Begin()
	if err := tx.Where("id = ?", id).Delete(&models.Gallery{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
