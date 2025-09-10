package event_dto

type CreateEventRequest struct {
	StartTime string `json:"start_time" validate:"required"`
	EndTime   string `json:"end_time"`
	Place     string `json:"place" validate:"required"`
	Address   string `json:"address" validate:"required"`
	MapsURL   string `json:"maps_url" validate:"required"`
}

type UpdateEventRequest struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Place     string `json:"place"`
	Address   string `json:"address"`
	MapsURL   string `json:"maps_url"`
}
