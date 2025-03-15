package models

import "time"

type Review struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	UserID            string          `gorm:"not null" json:"user_id"`
	User              *User           `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	InvitationThemeID uint            `gorm:"not null;index" json:"invitation_theme_id"`
	InvitationTheme   InvitationTheme `gorm:"foreignKey:InvitationThemeID;references:ID" json:"-"`
	Star              int             `gorm:"not null" json:"star"`
	Comment           string          `gorm:"type:text;not null" json:"comment"`
	CreatedAt         time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}
