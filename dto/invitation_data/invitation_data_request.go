package invitation_data_dto

type CreateInvitationDataRequest struct {
	EventName string `json:"event_name" binding:"required"`
	EventDate string `json:"event_date" binding:"required"` // Format ISO8601
	Location  string `json:"location" binding:"required"`
}

type UpdateInvitationDataRequest struct {
	EventName string `json:"event_name"`
	EventDate string `json:"event_date"` // Format ISO8601
	Location  string `json:"location"`
}
