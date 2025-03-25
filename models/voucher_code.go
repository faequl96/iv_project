package models

import "time"

type VoucherCode struct {
	ID                 uint   `gorm:"primaryKey" json:"id"`
	Name               string `gorm:"size:100;uniqueIndex;not null" json:"name"`
	DiscountPercentage uint   `json:"discount_percentage"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
