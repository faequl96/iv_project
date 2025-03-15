package models

import "time"

type InvitationData struct {
	ID           uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	InvitationID uint        `gorm:"not null;index" json:"invitation_id"`
	Invitation   *Invitation `gorm:"foreignKey:InvitationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	EventName    string      `gorm:"size:255;not null" json:"event_name"`
	EventDate    time.Time   `gorm:"not null" json:"event_date"`
	Location     string      `gorm:"size:255;not null" json:"location"`
	Gallery      []*Gallery  `gorm:"foreignKey:InvitationDataID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"gallery,omitempty"`
	CreatedAt    time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}
