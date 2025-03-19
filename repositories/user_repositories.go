package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type UserRepositories interface {
	CreateUser(user *models.User) error
	GetUserByID(id string) (*models.User, error)
	GetUsers() ([]models.User, error)
	DeleteUser(id string) error
}

func UserRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateUser(user *models.User) error {
	tx := r.db.Begin()
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *repository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("UserProfile").Preload("IVCoin").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("UserProfile").Preload("IVCoin").Find(&users).Error
	return users, err
}

func (r *repository) DeleteUser(id string) error {
	tx := r.db.Begin()
	if err := tx.Where("id = ?", id).Delete(&models.User{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
