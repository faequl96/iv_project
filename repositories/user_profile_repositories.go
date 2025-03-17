package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type UserProfileRepositories interface {
	GetUserProfileByID(id uint) (*models.UserProfile, error)
	UpdateUserProfile(userProfile *models.UserProfile) error
}

func UserProfileRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetUserProfileByID(id uint) (*models.UserProfile, error) {
	var userProfile models.UserProfile
	err := r.db.Where("id = ?", id).First(&userProfile).Error
	if err != nil {
		return nil, err
	}
	return &userProfile, nil
}

func (r *repository) UpdateUserProfile(userProfile *models.UserProfile) error {
	return r.db.Save(userProfile).Error
}
