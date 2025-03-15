package models

import "time"

type InvitationTheme struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	NormalPrice int       `gorm:"not null" json:"normal_price"`
	DiskonPrice int       `gorm:"not null" json:"diskon_price"`
	Category    string    `gorm:"type:varchar(100);not null" json:"category"`
	Reviews     []*Review `gorm:"foreignKey:InvitationThemeID" json:"reviews,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
