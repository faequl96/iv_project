package models

import "time"

type Review struct {
	ID                uint            `gorm:"primaryKey"`
	User              User            `gorm:"foreignKey:UserID"`
	InvitationThemeID uint            `gorm:"not null"`
	InvitationTheme   InvitationTheme `gorm:"foreignKey:InvitationThemeID"`
	Star              int             `gorm:"type:int;not null;check:star >= 1 AND star <= 5"`
	Comment           string          `gorm:"type:text;not null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
