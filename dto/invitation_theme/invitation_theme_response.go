package invitation_theme_dto

import (
	category_dto "iv_project/dto/category"
	review_dto "iv_project/dto/review"
)

type InvitationThemeResponse struct {
	ID          uint                            `json:"id"`
	Title       string                          `json:"title"`
	NormalPrice uint                            `json:"normal_price"`
	DiskonPrice uint                            `json:"diskon_price"`
	Categories  []category_dto.CategoryResponse `json:"categories,omitempty"`
	Reviews     []review_dto.ReviewResponse     `json:"reviews,omitempty"`
}
