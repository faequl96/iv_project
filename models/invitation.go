package models

import (
	"time"
)

type Invitation struct {
	ID                uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            uint            `gorm:"not null;index" json:"user_id"`
	Status            string          `gorm:"type:varchar(50);not null" json:"status"`
	InvitationThemeID uint            `gorm:"not null" json:"invitation_theme_id"`
	InvitationTheme   InvitationTheme `gorm:"foreignKey:InvitationThemeID" json:"invitation_theme"`
	InvitationData    InvitationData  `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;" json:"invitation_data"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}
