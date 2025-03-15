package models

import "time"

type Review struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	UserID            string          `gorm:"not null" json:"user_id"`
	User              *User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	InvitationThemeID uint            `gorm:"not null" json:"invitation_theme_id"`
	InvitationTheme   InvitationTheme `gorm:"foreignKey:InvitationThemeID" json:"-"`
	Star              int             `gorm:"not null" json:"star"`
	Comment           string          `gorm:"type:text;not null" json:"comment"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}
