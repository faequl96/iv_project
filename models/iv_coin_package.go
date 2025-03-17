package models

import "time"

type IVCoinPackage struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"size:100;not null;index" json:"name"`
	CoinAmount    uint      `gorm:"not null;default:0" json:"coin_amount"`
	Price         uint      `gorm:"not null;default:0" json:"price"`
	DiscountPrice uint      `gorm:"not null;default:0" json:"discount_price"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
