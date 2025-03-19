package invitation_theme_dto

type CreateInvitationThemeRequest struct {
	Title            string `json:"title" validate:"required"`
	IDRPrice         uint   `json:"idr_price" validate:"required"`
	IDRDiscountPrice uint   `json:"idr_discount_price" validate:"required"`
	IVCPrice         uint   `json:"ivc_price" validate:"required"`
	IVCDiscountPrice uint   `json:"ivc_discount_price" validate:"required"`
	Categories       []uint `json:"categories" validate:"required"`
}

type UpdateInvitationThemeRequest struct {
	Title            string `json:"title"`
	IDRPrice         uint   `json:"idr_price"`
	IDRDiscountPrice uint   `json:"idr_discount_price"`
	IVCPrice         uint   `json:"ivc_price"`
	IVCDiscountPrice uint   `json:"ivc_discount_price"`
	Categories       []uint `json:"categories"`
}
