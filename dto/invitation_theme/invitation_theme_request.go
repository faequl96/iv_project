package invitation_theme_dto

type CreateInvitationThemeRequest struct {
	Title       string `json:"title" binding:"required"`
	NormalPrice int    `json:"normal_price" binding:"required,gte=0"`
	DiskonPrice int    `json:"diskon_price" binding:"required,gte=0"`
	Category    string `json:"category" binding:"required"`
}

type UpdateInvitationThemeRequest struct {
	Title       string `json:"title"`
	NormalPrice int    `json:"normal_price"`
	DiskonPrice int    `json:"diskon_price"`
	Category    string `json:"category"`
}
