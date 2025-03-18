package models

import "time"

type Transaction struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductType     string         `gorm:"type:enum('invitation','ivcoin');not null" json:"product_type"`
	Status          string         `gorm:"type:enum('pending','completed');not null;default:'pending'" json:"status"`
	PaymentMethod   string         `gorm:"type:enum('transfer','ivcoin');not null" json:"payment_method"`
	ReferenceNumber string         `gorm:"size:100;uniqueIndex;not null" json:"reference_number"`
	IDRPrice        uint           `gorm:"not null" json:"idr_price"`
	IDRDiscount     uint           `gorm:"not null" json:"idr_discount"`
	IDRTotalPrice   uint           `gorm:"not null" json:"idr_total_price"`
	IVCPrice        uint           `gorm:"not null" json:"ivc_price"`
	IVCDiscount     uint           `gorm:"not null" json:"ivc_discount"`
	IVCTotalPrice   uint           `gorm:"not null" json:"ivc_total_price"`
	ProductID       uint           `gorm:"index" json:"product_id,omitempty"`
	Invitation      *Invitation    `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"invitation,omitempty"`
	IVCoinPackage   *IVCoinPackage `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"iv_coin_package,omitempty"`

	UserID string `gorm:"size:36;not null;index" json:"user_id"`
	User   *User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
