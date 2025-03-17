package models

type Category struct {
	ID     uint              `gorm:"primaryKey" json:"id"`
	Name   string            `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Themes []InvitationTheme `gorm:"many2many:invitation_theme_categories;" json:"invitation_themes,omitempty"`
}
