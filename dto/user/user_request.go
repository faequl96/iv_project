package user_dto

import user_profile_dto "iv_project/dto/user_profile"

type CreateUserRequest struct {
	ID          string                                     `json:"id" binding:"required"` // Firebase UID
	UserProfile *user_profile_dto.CreateUserProfileRequest `json:"user_profile" binding:"required"`
}
