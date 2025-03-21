package discount_dto

type DiscountRequest struct {
	DiscountCategory string `json:"discount_category" validate:"required"`
	Percentage       uint   `json:"percentage" validate:"required"`
}
