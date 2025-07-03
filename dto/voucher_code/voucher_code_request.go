package voucher_code_dto

type VoucherCodeRequest struct {
	Name               string   `json:"name" validate:"required"`
	DiscountPercentage uint     `json:"discount_percentage" validate:"required"`
	UsageLimitPerUser  int      `json:"usage_limit_per_user" validate:"required"`
	IsGlobal           bool     `json:"is_global"`
	UserIDs            []string `json:"user_ids,omitempty"`
}
