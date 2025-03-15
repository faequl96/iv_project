package models

import (
	"time"

	"gorm.io/gorm"
)

type Invitation struct {
	ID                uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            string          `gorm:"not null;index" json:"user_id"`
	User              *User           `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Status            string          `gorm:"type:enum('pending','approved','rejected');not null;default:'pending'" json:"status"`
	InvitationThemeID uint            `gorm:"not null;index" json:"invitation_theme_id"`
	InvitationTheme   InvitationTheme `gorm:"foreignKey:InvitationThemeID;references:ID" json:"invitation_theme"`
	InvitationData    *InvitationData `gorm:"foreignKey:InvitationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"invitation_data,omitempty"`
	CreatedAt         time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt  `gorm:"index" json:"deleted_at"`
}
