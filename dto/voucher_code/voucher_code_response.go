package voucher_code_dto

type VoucherCodeResponse struct {
	ID                 uint   `json:"id"`
	Name               string `json:"name"`
	DiscountPercentage uint   `json:"discount_percentage"`
}
