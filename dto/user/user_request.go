package user_dto

import "iv_project/models"

type UpdateUserRequest struct {
	Role models.UserRoleType `json:"role" validate:"required"`
}
