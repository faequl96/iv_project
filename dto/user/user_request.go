package user_dto

type CreateUserRequest struct {
	ID string `json:"id" validate:"required"` // Firebase UID
}
