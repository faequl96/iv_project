package invitation_dto

type InvitationRequest struct {
	UserID            string                `json:"user_id" binding:"required"`
	InvitationThemeID uint                  `json:"invitation_theme_id" binding:"required"`
	Status            string                `json:"status" binding:"required"`
	InvitationData    InvitationDataRequest `json:"invitation_data"`
}

type InvitationDataRequest struct {
	EventName string         `json:"event_name" binding:"required"`
	EventDate string         `json:"event_date" binding:"required"` // Format ISO8601
	Location  string         `json:"location" binding:"required"`
	Gallery   GalleryRequest `json:"gallery"`
}

type GalleryRequest struct {
	ImageURL1  string `json:"image_url_1"`
	ImageURL2  string `json:"image_url_2"`
	ImageURL3  string `json:"image_url_3"`
	ImageURL4  string `json:"image_url_4"`
	ImageURL5  string `json:"image_url_5"`
	ImageURL6  string `json:"image_url_6"`
	ImageURL7  string `json:"image_url_7"`
	ImageURL8  string `json:"image_url_8"`
	ImageURL9  string `json:"image_url_9"`
	ImageURL10 string `json:"image_url_10"`
	ImageURL11 string `json:"image_url_11"`
	ImageURL12 string `json:"image_url_12"`
}
