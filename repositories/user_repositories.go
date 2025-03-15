package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type UserRepositories interface {
	CreateUser(user models.User) error
	GetUserByID(userID string) (models.User, error)
	// UpdateUser(user models.User) (models.User, error)
}

func UserRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateUser(user models.User) error {
	return r.db.Create(&user).Error
}

func (r *repository) GetUserByID(userID string) (models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", userID).First(&user).Error
	return user, err
}

// func (r *repository) UpdateUser(user models.User) (models.User, error) {
// 	err := r.db.Save(&user).Error
// 	return user, err
// }
