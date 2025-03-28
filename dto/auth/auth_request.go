package auth_dto

type AuthRequest struct {
	ID    string `json:"id" validate:"required"` // Firebase UID
	Email string `json:"email" validate:"required"`
}
