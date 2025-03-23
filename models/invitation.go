package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type InvitationStatusType string

const (
	InvitationStatusDraft  InvitationStatusType = "draft"
	InvitationStatusActive InvitationStatusType = "active"
)

func (u InvitationStatusType) String() string {
	maps := map[InvitationStatusType]string{
		InvitationStatusDraft:  "draft",
		InvitationStatusActive: "active",
	}
	return maps[u]
}

func StringToInvitationStatusType(value string) InvitationStatusType {
	maps := map[string]InvitationStatusType{
		"draft":  InvitationStatusDraft,
		"active": InvitationStatusActive,
	}
	return maps[value]
}

type Invitation struct {
	ID                  uint                 `gorm:"primaryKey;autoIncrement" json:"id"`
	Status              InvitationStatusType `gorm:"type:varchar(50);not null;default:'draft'" json:"status"`
	InvitationThemeID   uint                 `gorm:"not null;index" json:"invitation_theme_id"`
	InvitationThemeName string               `gorm:"size:150;not null;index" json:"invitation_theme_name"`
	InvitationData      *InvitationData      `gorm:"foreignKey:InvitationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"invitation_data,omitempty"`

	InvitationTheme *InvitationTheme `gorm:"foreignKey:InvitationThemeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"invitation_theme,omitempty"`
	UserID          string           `gorm:"size:36;not null;index" json:"user_id"`
	User            *User            `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time             `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time             `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index" json:"-"`
}
