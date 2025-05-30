package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type UserProfileRepositories interface {
	GetUserProfileByID(id uint) (*models.UserProfile, error)
	GetUserProfileByUserID(userId string) (*models.UserProfile, error)
	UpdateUserProfile(userProfile *models.UserProfile) error
}

func UserProfileRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetUserProfileByID(id uint) (*models.UserProfile, error) {
	var userProfile models.UserProfile
	err := r.db.First(&userProfile, id).Error
	if err != nil {
		return nil, err
	}
	return &userProfile, nil
}

func (r *repository) GetUserProfileByUserID(userID string) (*models.UserProfile, error) {
	var userProfile models.UserProfile
	err := r.db.First(&userProfile, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &userProfile, nil
}

func (r *repository) UpdateUserProfile(userProfile *models.UserProfile) error {
	return r.db.Save(userProfile).Error
}
