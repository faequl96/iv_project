package invitation_theme_dto

type CreateInvitationThemeRequest struct {
	Title            string `json:"title" binding:"required"`
	IDRPrice         uint   `json:"idr_price" binding:"required,gte=0"`
	IDRDiscountPrice uint   `json:"idr_discount_price" binding:"required,gte=0"`
	IVCPrice         uint   `json:"ivc_price" binding:"required,gte=0"`
	IVCDiscountPrice uint   `json:"ivc_discount_price" binding:"required,gte=0"`
	Categories       []uint `json:"categories" binding:"required"`
}

type UpdateInvitationThemeRequest struct {
	Title            string `json:"title"`
	IDRPrice         uint   `json:"idr_price"`
	IDRDiscountPrice uint   `json:"idr_discount_price"`
	IVCPrice         uint   `json:"ivc_price"`
	IVCDiscountPrice uint   `json:"ivc_discount_price"`
	Categories       []uint `json:"categories"`
}
