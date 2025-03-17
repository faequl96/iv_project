package user_dto

import (
	iv_coin_dto "iv_project/dto/iv_coin"
	user_profile_dto "iv_project/dto/user_profile"
)

type UserResponse struct {
	ID          string                                `json:"id"` // Firebase UID
	UserProfile *user_profile_dto.UserProfileResponse `json:"user_profile"`
	IVCoin      *iv_coin_dto.IVCoinResponse           `json:"iv_coin"`
}
