package user_dto

type UserRequest struct {
	Role string `json:"role" validate:"required"`
}
