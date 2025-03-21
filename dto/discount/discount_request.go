package discount_dto

type DiscountRequest struct {
	DiscountCategoryID uint `json:"discount_category_id" validate:"required"`
	Percentage         uint `json:"percentage"`
}
