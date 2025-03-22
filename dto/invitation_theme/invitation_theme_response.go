package invitation_theme_dto

import (
	category_dto "iv_project/dto/category"
	discount_category_dto "iv_project/dto/discount_category"
	review_dto "iv_project/dto/review"
)

type InvitationThemeResponse struct {
	ID                 uint                                             `json:"id"`
	Name               string                                           `json:"name"`
	IDRPrice           uint                                             `json:"idr_price"`
	IDRDiscountPrice   uint                                             `json:"idr_discount_price"`
	IVCPrice           uint                                             `json:"ivc_price"`
	IVCDiscountPrice   uint                                             `json:"ivc_discount_price"`
	Categories         []category_dto.CategoryResponse                  `json:"categories"`
	DiscountCategories []discount_category_dto.DiscountCategoryResponse `json:"discount_categories"`
	Reviews            []review_dto.ReviewResponse                      `json:"reviews"`
}
