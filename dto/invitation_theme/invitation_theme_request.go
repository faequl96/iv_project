package invitation_dto

type InvitationThemeRequest struct {
	Title       string  `json:"title" binding:"required"`
	NormalPrice int     `json:"normal_price" binding:"required,gte=0"`
	DiskonPrice int     `json:"diskon_price" binding:"required,gte=0"`
	Category    string  `json:"category" binding:"required"`
	Rating      float32 `json:"rating" binding:"required,gte=0,lte=5"`
}
