package auth_dto

import user_dto "iv_project/dto/user"

type AuthResponse struct {
	Token string                `json:"token"`
	User  user_dto.UserResponse `json:"user"`
}
