package models

import "time"

type IVCoinPackage struct {
	ID                 uint               `gorm:"primaryKey;autoIncrement" json:"id"`
	Name               string             `gorm:"size:100;not null;index" json:"name"`
	CoinAmount         uint               `gorm:"not null;default:0" json:"coin_amount"`
	IDRPrice           uint               `gorm:"not null;default:0" json:"idr_price"`
	IDRDiscountPrice   uint               `gorm:"not null;default:0" json:"idr_discount_price"`
	IVCPrice           uint               `gorm:"not null;default:0" json:"ivc_price"`
	IVCDiscountPrice   uint               `gorm:"not null;default:0" json:"ivc_discount_price"`
	DiscountCategories []DiscountCategory `gorm:"many2many:iv_coin_package_discount_categories;" json:"discount_categories,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
