package discount_category_dto

type CreateDiscountCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateDiscountCategoryRequest struct {
	Name string `json:"name"`
}
