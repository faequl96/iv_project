package auth_dto

type LoginAuthRequest struct {
	ID    string `json:"id" validate:"required"` // Firebase UID
	Email string `json:"email" validate:"required"`
}
