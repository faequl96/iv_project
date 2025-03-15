package models

import (
	"time"
)

type Invitation struct {
	ID                uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            string          `gorm:"not null;index" json:"user_id"`
	User              User            `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Status            string          `gorm:"type:enum('pending','approved','rejected');default:'pending';not null" json:"status"`
	InvitationThemeID uint            `gorm:"not null" json:"invitation_theme_id"`
	InvitationTheme   InvitationTheme `gorm:"foreignKey:InvitationThemeID" json:"invitation_theme"`
	InvitationData    *InvitationData `gorm:"foreignKey:InvitationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"invitation_data,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}
