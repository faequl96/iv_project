package models

import "time"

type InvitationData struct {
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID           uint      `gorm:"not null;uniqueIndex;constraint:OnDelete:CASCADE;" json:"order_id"`
	EventName         string    `gorm:"type:varchar(255);not null" json:"event_name"`
	EventDate         time.Time `gorm:"not null" json:"event_date"`
	Location          string    `gorm:"type:varchar(255);not null" json:"location"`
	GalleryImageURL1  string    `gorm:"type:varchar(255)" json:"gallery_image_url_1"`
	GalleryImageURL2  string    `gorm:"type:varchar(255)" json:"gallery_image_url_2"`
	GalleryImageURL3  string    `gorm:"type:varchar(255)" json:"gallery_image_url_3"`
	GalleryImageURL4  string    `gorm:"type:varchar(255)" json:"gallery_image_url_4"`
	GalleryImageURL5  string    `gorm:"type:varchar(255)" json:"gallery_image_url_5"`
	GalleryImageURL6  string    `gorm:"type:varchar(255)" json:"gallery_image_url_6"`
	GalleryImageURL7  string    `gorm:"type:varchar(255)" json:"gallery_image_url_7"`
	GalleryImageURL8  string    `gorm:"type:varchar(255)" json:"gallery_image_url_8"`
	GalleryImageURL9  string    `gorm:"type:varchar(255)" json:"gallery_image_url_9"`
	GalleryImageURL10 string    `gorm:"type:varchar(255)" json:"gallery_image_url_10"`
	GalleryImageURL11 string    `gorm:"type:varchar(255)" json:"gallery_image_url_11"`
	GalleryImageURL12 string    `gorm:"type:varchar(255)" json:"gallery_image_url_12"`
}
