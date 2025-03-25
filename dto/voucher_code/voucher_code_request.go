package voucher_code_dto

type CreateVoucherCodeRequest struct {
	Name               string `json:"name" validate:"required"`
	DiscountPercentage uint   `json:"discount_percentage" validate:"required"`
}

type UpdateVoucherCodeRequest struct {
	Name               string `json:"name" validate:"required"`
	DiscountPercentage uint   `json:"discount_percentage" validate:"required"`
}
