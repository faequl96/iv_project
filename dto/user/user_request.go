package user_dto

import "iv_project/models"

type CreateUserRequest struct {
	ID string `json:"id" validate:"required"` // Firebase UID
}

type UpdateUserRequest struct {
	Role models.UserRoleType `json:"role" validate:"required"`
}
