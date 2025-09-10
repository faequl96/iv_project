package event_dto

type EventResponse struct {
	ID        uint   `json:"id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Place     string `json:"place"`
	Address   string `json:"address"`
	MapsURL   string `json:"maps_url"`
}
