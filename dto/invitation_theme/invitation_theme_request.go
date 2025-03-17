package invitation_theme_dto

type CreateInvitationThemeRequest struct {
	Title         string `json:"title" binding:"required"`
	Price         uint   `json:"price" binding:"required,gte=0"`
	DiscountPrice uint   `json:"discount_price" binding:"required,gte=0"`
	Categories    []uint `json:"categories" binding:"required"`
}

type UpdateInvitationThemeRequest struct {
	Title         string `json:"title"`
	Price         uint   `json:"price"`
	DiscountPrice uint   `json:"discount_price"`
	Categories    []uint `json:"categories"`
}
