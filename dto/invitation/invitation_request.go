package invitation_dto

type InvitationRequest struct {
	UserID         string                `json:"user_id" binding:"required"`
	Status         string                `json:"status" binding:"required"`
	InvitationData InvitationDataRequest `json:"invitation_data" binding:"required"`
}

type InvitationDataRequest struct {
	EventName         string `json:"event_name" binding:"required"`
	EventDate         string `json:"event_date" binding:"required"`
	Location          string `json:"location" binding:"required"`
	GalleryImageURL1  string `json:"gallery_image_url_1"`
	GalleryImageURL2  string `json:"gallery_image_url_2"`
	GalleryImageURL3  string `json:"gallery_image_url_3"`
	GalleryImageURL4  string `json:"gallery_image_url_4"`
	GalleryImageURL5  string `json:"gallery_image_url_5"`
	GalleryImageURL6  string `json:"gallery_image_url_6"`
	GalleryImageURL7  string `json:"gallery_image_url_7"`
	GalleryImageURL8  string `json:"gallery_image_url_8"`
	GalleryImageURL9  string `json:"gallery_image_url_9"`
	GalleryImageURL10 string `json:"gallery_image_url_10"`
	GalleryImageURL11 string `json:"gallery_image_url_11"`
	GalleryImageURL12 string `json:"gallery_image_url_12"`
}
