package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type InvitationTheme struct {
	ID            uint                  `gorm:"primaryKey" json:"id"`
	Title         string                `gorm:"size:255;not null;index" json:"title"`
	Price         uint                  `gorm:"not null" json:"price"`
	DiscountPrice uint                  `gorm:"not null" json:"discount_price"`
	Categories    []Category            `gorm:"many2many:invitation_theme_categories;" json:"categories,omitempty"`
	Reviews       []Review              `gorm:"foreignKey:InvitationThemeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"reviews,omitempty"`
	CreatedAt     time.Time             `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time             `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     soft_delete.DeletedAt `gorm:"index" json:"-"`
}
