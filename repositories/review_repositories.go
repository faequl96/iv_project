package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type ReviewRepositories interface {
	CreateReview(review models.Review) error
	GetReviewByID(reviewID uint) (models.Review, error)
	GetReviewsByInvitationThemeID(invitationThemeID uint) ([]models.Review, error)
	UpdateReview(review models.Review) error
	DeleteReview(review models.Review) error
}

func ReviewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateReview(review models.Review) error {
	return r.db.Create(&review).Error
}

func (r *repository) GetReviewByID(reviewID uint) (models.Review, error) {
	var review models.Review
	err := r.db.Where("id = ?", reviewID).First(&review).Error
	return review, err
}

func (r *repository) GetReviewsByInvitationThemeID(invitationThemeID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("id = ?", invitationThemeID).Find(&reviews).Error
	return reviews, err
}

func (r *repository) UpdateReview(review models.Review) error {
	err := r.db.Save(&review).Error
	return err
}

func (r *repository) DeleteReview(review models.Review) error {
	err := r.db.Delete(&review).Error
	return err
}
