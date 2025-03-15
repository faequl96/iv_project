package models

import (
	"time"
)

type Invitation struct {
	ID                uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            string          `gorm:"not null;index" json:"user_id"`
	User              *User           `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Status            string          `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	InvitationThemeID uint            `gorm:"not null" json:"invitation_theme_id"`
	InvitationTheme   InvitationTheme `gorm:"foreignKey:InvitationThemeID" json:"invitation_theme"`
	InvitationData    *InvitationData `gorm:"foreignKey:InvitationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"invitation_data,omitempty"`
	CreatedAt         time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}
