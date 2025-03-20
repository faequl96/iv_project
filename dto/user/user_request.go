package user_dto

type UpdateUserRequest struct {
	Role string `json:"role" validate:"required"`
}
