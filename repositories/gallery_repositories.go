package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type GalleryRepositories interface {
	GetGalleryByID(galleryID uint) (*models.Gallery, error)
	DeleteGallery(gallery *models.Gallery) error
}

func GalleryRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetGalleryByID(galleryID uint) (*models.Gallery, error) {
	var gallery models.Gallery
	err := r.db.Where("id = ?", galleryID).First(&gallery).Error
	if err != nil {
		return nil, err
	}
	return &gallery, nil
}

func (r *repository) DeleteGallery(gallery *models.Gallery) error {
	return r.db.Delete(gallery).Error
}
