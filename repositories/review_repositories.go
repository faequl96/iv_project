package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type ReviewRepositories interface {
	CreateReview(review models.Review) error
	GetReviewByID(reviewID uint) (models.Review, error)
	GetReviews() ([]models.Review, error)
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

func (r *repository) GetReviews() ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Find(&reviews).Error
	return reviews, err
}

func (r *repository) GetReviewsByInvitationThemeID(invitationThemeID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("invitation_theme_id = ?", invitationThemeID).Find(&reviews).Error
	return reviews, err
}

func (r *repository) UpdateReview(review models.Review) error {
	return r.db.Save(&review).Error
}

func (r *repository) DeleteReview(review models.Review) error {
	return r.db.Delete(&review).Error
}
