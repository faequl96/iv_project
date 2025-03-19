package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type ReviewRepositories interface {
	CreateReview(review *models.Review) error
	GetReviewByID(id uint) (*models.Review, error)
	GetReviewsByInvitationThemeID(invitationThemeID uint) ([]models.Review, error)
	UpdateReview(review *models.Review) error
	DeleteReview(id uint) error
}

func ReviewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateReview(review *models.Review) error {
	return r.db.Create(review).Error
}

func (r *repository) GetReviewByID(id uint) (*models.Review, error) {
	var review models.Review
	err := r.db.First(&review, id).Error
	return &review, err
}

func (r *repository) GetReviewsByInvitationThemeID(invitationThemeID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("invitation_theme_id = ?", invitationThemeID).Find(&reviews).Error
	return reviews, err
}

func (r *repository) UpdateReview(review *models.Review) error {
	return r.db.Model(review).Updates(review).Error
}

func (r *repository) DeleteReview(id uint) error {
	return r.db.Delete(&models.Review{}, id).Error
}
