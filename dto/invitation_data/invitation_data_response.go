package invitation_data_dto

import gallery_dto "iv_project/dto/gallery"

type InvitationDataResponse struct {
	ID           uint                         `json:"id"`
	EventName    string                       `json:"event_name"`
	EventDate    string                       `json:"event_date"` // Format ISO8601
	Location     string                       `json:"location"`
	MainImageURL string                       `json:"main_image_url"`
	Gallery      *gallery_dto.GalleryResponse `json:"gallery"`
}
