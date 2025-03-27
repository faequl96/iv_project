package voucher_code_dto

type VoucherCodeRequest struct {
	Name               string `json:"name" validate:"required"`
	DiscountPercentage uint   `json:"discount_percentage" validate:"required"`
}
