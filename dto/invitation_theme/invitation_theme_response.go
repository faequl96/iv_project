package invitation_theme_dto

import (
	category_dto "iv_project/dto/category"
	review_dto "iv_project/dto/review"
)

type InvitationThemeResponse struct {
	ID               uint                            `json:"id"`
	Title            string                          `json:"title"`
	IDRPrice         uint                            `json:"idr_price"`
	IDRDiscountPrice uint                            `json:"idr_discount_price"`
	IVCPrice         uint                            `json:"ivc_price"`
	IVCDiscountPrice uint                            `json:"ivc_discount_price"`
	Categories       []category_dto.CategoryResponse `json:"categories"`
	Reviews          []review_dto.ReviewResponse     `json:"reviews"`
}
