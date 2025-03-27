package discount_category_dto

type DiscountCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}
