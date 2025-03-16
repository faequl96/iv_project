package invitation_dto

type InvitationRequest struct {
	UserID            string                `json:"user_id" binding:"required"`
	InvitationThemeID uint                  `json:"invitation_theme_id" binding:"required"`
	Status            string                `json:"status" binding:"required"`
	InvitationData    InvitationDataRequest `json:"invitation_data"`
}

type InvitationDataRequest struct {
	EventName string    `json:"event_name" binding:"required"`
	EventDate string    `json:"event_date" binding:"required"` // Format ISO8601
	Location  string    `json:"location" binding:"required"`
	Gallery   []Gallery `json:"gallery"`
}

type Gallery struct {
	Position int    `json:"position" binding:"required"`
	ImageURL string `json:"image_url" binding:"required"`
}
