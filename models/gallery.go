package models

type Gallery struct {
	ID               uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	InvitationDataID uint   `gorm:"not null" json:"invitation_data_id"`
	ImageURL         string `gorm:"type:varchar(255);not null" json:"image_url"`
}
