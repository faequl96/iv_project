package invitation_theme_dto

type CreateInvitationThemeRequest struct {
	Title       string `json:"title" binding:"required"`
	NormalPrice uint   `json:"normal_price" binding:"required,gte=0"`
	DiskonPrice uint   `json:"diskon_price" binding:"required,gte=0"`
	Category    string `json:"category" binding:"required"`
}

type UpdateInvitationThemeRequest struct {
	Title       string `json:"title"`
	NormalPrice uint   `json:"normal_price"`
	DiskonPrice uint   `json:"diskon_price"`
	Category    string `json:"category"`
}
