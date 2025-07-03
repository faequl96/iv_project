package voucher_code_dto

import (
	user_dto "iv_project/dto/user"
)

type VoucherCodeResponse struct {
	ID                 uint                    `json:"id"`
	Name               string                  `json:"name"`
	DiscountPercentage uint                    `json:"discount_percentage"`
	UsageLimitPerUser  int                     `json:"usage_limit_per_user"`
	IsGlobal           bool                    `json:"is_global"`
	Users              []user_dto.UserResponse `json:"users"`
}
