package models

type Gallery struct {
	ID               uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	InvitationDataID uint            `gorm:"not null;index" json:"invitation_data_id"`
	InvitationData   *InvitationData `gorm:"foreignKey:InvitationDataID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Position         int             `gorm:"not null" json:"position"`
	ImageURL         string          `gorm:"size:255;not null" json:"image_url"`
}
