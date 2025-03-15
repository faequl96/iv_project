package models

import "time"

type InvitationData struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	InvitationID uint       `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"invitation_id"`
	Invitation   Invitation `gorm:"foreignKey:InvitationID" json:"-"`
	EventName    string     `gorm:"type:varchar(255);not null" json:"event_name"`
	EventDate    time.Time  `gorm:"not null" json:"event_date"`
	Location     string     `gorm:"type:varchar(255);not null" json:"location"`
	Gallery      []Gallery  `gorm:"foreignKey:InvitationDataID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"gallery"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
