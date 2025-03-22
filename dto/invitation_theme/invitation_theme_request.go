package invitation_theme_dto

type CreateInvitationThemeRequest struct {
	Name                string `json:"name" validate:"required"`
	IDRPrice            uint   `json:"idr_price" validate:"required"`
	IDRDiscountPrice    uint   `json:"idr_discount_price" validate:"required"`
	IVCPrice            uint   `json:"ivc_price" validate:"required"`
	IVCDiscountPrice    uint   `json:"ivc_discount_price" validate:"required"`
	CategoryIDs         []uint `json:"category_ids" validate:"required"`
	DiscountCategoryIDs []uint `json:"discount_category_ids" validate:"required"`
}

type UpdateInvitationThemeRequest struct {
	Name                string `json:"n"`
	IDRPrice            uint   `json:"idr_price"`
	IDRDiscountPrice    uint   `json:"idr_discount_price"`
	IVCPrice            uint   `json:"ivc_price"`
	IVCDiscountPrice    uint   `json:"ivc_discount_price"`
	CategoryIDs         []uint `json:"category_ids"`
	DiscountCategoryIDs []uint `json:"discount_category_ids"`
}
