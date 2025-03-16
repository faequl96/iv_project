package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type UserRepositories interface {
	CreateUser(user models.User) error
	GetUserByID(userID string) (models.User, error)
	GetUsers() ([]models.User, error)
	UpdateUser(user models.User) error
	DeleteUser(user models.User) error
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

func (r *repository) GetUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *repository) UpdateUser(user models.User) error {
	return r.db.Save(&user).Error
}

func (r *repository) DeleteUser(user models.User) error {
	return r.db.Delete(&user).Error
}
