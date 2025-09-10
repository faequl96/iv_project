package handlers

import (
	event_dto "iv_project/dto/event"
	"iv_project/models"
	"time"
)

func ConvertToEventResponse(event models.Event) event_dto.EventResponse {
	return event_dto.EventResponse{
		ID:        event.ID,
		StartTime: event.StartTime.Format(time.RFC3339),
		EndTime:   event.EndTime.Format(time.RFC3339),
		Place:     event.Place,
		Address:   event.Address,
		MapsURL:   event.MapsURL,
	}
}
