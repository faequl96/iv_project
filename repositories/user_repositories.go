package repositories

import (
	query_dto "iv_project/dto/query"
	"iv_project/models"
	"iv_project/pkg/utils"

	"gorm.io/gorm"
)

type UserRepositories interface {
	CreateUser(user *models.User) error
	GetUserByID(id string) (*models.User, error)
	GetUsers(query *query_dto.QueryRequest) ([]models.User, error)
	GetUsersByIDs(ids []string) ([]models.User, error)
	UpdateUser(user *models.User) error
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

func (r *repository) GetUsers(query *query_dto.QueryRequest) ([]models.User, error) {
	var users []models.User
	db := utils.ImplementQuery(r.db, query)
	err := db.Preload("UserProfile").Preload("IVCoin").Find(&users).Error
	return users, err
}

func (r *repository) GetUsersByIDs(ids []string) ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users, "id IN ?", ids).Error
	return users, err
}

func (r *repository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *repository) DeleteUser(id string) error {
	tx := r.db.Begin()
	if err := tx.Delete(&models.User{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
