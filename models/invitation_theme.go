package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type InvitationTheme struct {
	ID          uint                  `gorm:"primaryKey" json:"id"`
	Title       string                `gorm:"type:varchar(255);not null" json:"title"`
	NormalPrice int                   `gorm:"not null" json:"normal_price"`
	DiskonPrice int                   `gorm:"not null" json:"diskon_price"`
	Category    string                `gorm:"type:varchar(100);not null" json:"category"`
	Reviews     []Review              `gorm:"foreignKey:InvitationThemeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"reviews,omitempty"`
	CreatedAt   time.Time             `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time             `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   soft_delete.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
