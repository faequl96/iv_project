package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type InvitationTheme struct {
	ID                 uint               `gorm:"primaryKey" json:"id"`
	Title              string             `gorm:"size:255;not null;index" json:"title"`
	IDRPrice           uint               `gorm:"not null;default:0" json:"idr_price"`
	IDRDiscountPrice   uint               `gorm:"not null;default:0" json:"idr_discount_price"`
	IVCPrice           uint               `gorm:"not null;default:0" json:"ivc_price"`
	IVCDiscountPrice   uint               `gorm:"not null;default:0" json:"ivc_discount_price"`
	Categories         []Category         `gorm:"many2many:invitation_theme_categories;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"categories,omitempty"`
	DiscountCategories []DiscountCategory `gorm:"many2many:invitation_theme_discount_categories;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"discount_categories,omitempty"`
	Reviews            []Review           `gorm:"foreignKey:InvitationThemeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"reviews,omitempty"`

	CreatedAt time.Time             `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time             `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index" json:"-"`
}
