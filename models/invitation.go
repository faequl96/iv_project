package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Invitation struct {
	ID                uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	Status            string           `gorm:"type:enum('pending','approved','rejected');not null;default:'pending'" json:"status"`
	InvitationData    *InvitationData  `gorm:"foreignKey:InvitationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"invitation_data,omitempty"`
	InvitationThemeID uint             `gorm:"not null;index" json:"invitation_theme_id"`
	InvitationTheme   *InvitationTheme `gorm:"foreignKey:InvitationThemeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"invitation_theme,omitempty"`

	UserID string `gorm:"size:36;not null;index" json:"user_id"`
	User   *User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time             `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time             `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index" json:"-"`
}
