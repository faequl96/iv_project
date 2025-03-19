package user_profile_dto

type CreateUserProfileRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`

	UserID string `json:"user_id" validate:"required"`
}

type UpdateUserProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
