package models

import "time"

type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:100;uniqueIndex;not null" json:"name"`

	InvitationThemes []InvitationTheme `gorm:"many2many:invitation_theme_categories;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
