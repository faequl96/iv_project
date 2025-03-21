package models

import "time"

type DiscountCategory struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:100;uniqueIndex;not null" json:"name"`

	InvitationThemes []InvitationTheme `gorm:"many2many:invitation_theme_discount_categories;" json:"-"`
	IVCoinPackage    []IVCoinPackage   `gorm:"many2many:iv_coin_package_discount_categories;" json:"-"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
