package user_dto

import (
	iv_coin_dto "iv_project/dto/iv_coin"
)

// UserResponse defines the response structure for user-related API responses
type UserResponse struct {
	ID       string                     `json:"id"`        // Firebase UID, uniquely identifies the user
	Email    string                     `json:"email"`     // User's email address
	UserName string                     `json:"user_name"` // Username of the user
	FullName string                     `json:"full_name"` // Full name of the user
	IVCoin   iv_coin_dto.IVCoinResponse `json:"iv_coin"`   // Associated IVCoin balance details
}
