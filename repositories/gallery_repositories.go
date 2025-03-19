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
	err := r.db.First(&Gallery, id).Error
	return &Gallery, err
}

func (r *repository) DeleteGallery(id uint) error {
	return r.db.Delete(&models.Gallery{}, id).Error
}
