package invitation_dto

type InvitationRequest struct {
	UserID         uint                  `json:"user_id" binding:"required"`
	Status         string                `json:"status" binding:"required"`
	InvitationData InvitationDataRequest `json:"invitation_data" binding:"required"`
}

type InvitationDataRequest struct {
	EventName string `json:"event_name" binding:"required"`
	EventDate string `json:"event_date" binding:"required"`
	Location  string `json:"location" binding:"required"`
}
