package invitation_data_dto

type CreateInvitationDataRequest struct {
	EventName string `json:"event_name" validate:"required"`
	EventDate string `json:"event_date" validate:"required"` // Format ISO8601
	Location  string `json:"location" validate:"required"`
}

type UpdateInvitationDataRequest struct {
	EventName string `json:"event_name"`
	EventDate string `json:"event_date"` // Format ISO8601
	Location  string `json:"location"`
}
