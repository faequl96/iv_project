package repositories

import (
	"iv_project/models"

	"gorm.io/gorm"
)

type UserRepositories interface {
	CreateUser(user *models.User) error
	GetUserByID(id string) (*models.User, error)
	GetUsers() ([]*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id string) error
}

func UserRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *repository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("IVCoin").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetUsers() ([]*models.User, error) {
	var users []*models.User
	err := r.db.Preload("IVCoin").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *repository) DeleteUser(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.User{}).Error
}
