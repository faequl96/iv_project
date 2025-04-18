package invitation_theme_dto

type CreateInvitationThemeRequest struct {
	Name                string `json:"name" validate:"required"`
	IDRPrice            uint   `json:"idr_price" validate:"required"`
	IVCPrice            uint   `json:"ivc_price" validate:"required"`
	CategoryIDs         []uint `json:"category_ids" validate:"required"`
	DiscountCategoryIDs []uint `json:"discount_category_ids" validate:"required"`
}

type UpdateInvitationThemeRequest struct {
	Name                string `json:"name"`
	IDRPrice            uint   `json:"idr_price"`
	IVCPrice            uint   `json:"ivc_price"`
	CategoryIDs         []uint `json:"category_ids"`
	DiscountCategoryIDs []uint `json:"discount_category_ids"`
}
