package models

import "time"

type Review struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Star    int    `gorm:"not null" json:"star"`
	Comment string `gorm:"type:text;not null" json:"comment"`
	UserID  string `gorm:"size:36;not null;index" json:"user_id"`
	User    *User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`

	InvitationThemeID uint             `gorm:"not null;index" json:"invitation_theme_id"`
	InvitationTheme   *InvitationTheme `gorm:"foreignKey:InvitationThemeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
