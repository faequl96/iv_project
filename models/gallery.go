package models

type Gallery struct {
	ID               uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	InvitationDataID uint            `gorm:"not null;index" json:"invitation_data_id"`
	InvitationData   *InvitationData `gorm:"foreignKey:InvitationDataID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	ImageURL1        string          `gorm:"type:varchar(255)" json:"image_url_1"`
	ImageURL2        string          `gorm:"type:varchar(255)" json:"image_url_2"`
	ImageURL3        string          `gorm:"type:varchar(255)" json:"image_url_3"`
	ImageURL4        string          `gorm:"type:varchar(255)" json:"image_url_4"`
	ImageURL5        string          `gorm:"type:varchar(255)" json:"image_url_5"`
	ImageURL6        string          `gorm:"type:varchar(255)" json:"image_url_6"`
	ImageURL7        string          `gorm:"type:varchar(255)" json:"image_url_7"`
	ImageURL8        string          `gorm:"type:varchar(255)" json:"image_url_8"`
	ImageURL9        string          `gorm:"type:varchar(255)" json:"image_url_9"`
	ImageURL10       string          `gorm:"type:varchar(255)" json:"image_url_10"`
	ImageURL11       string          `gorm:"type:varchar(255)" json:"image_url_11"`
	ImageURL12       string          `gorm:"type:varchar(255)" json:"image_url_12"`
}
