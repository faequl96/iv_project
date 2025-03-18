package models

import "time"

type ProductType string

const (
	ProductInvitation    ProductType = "invitation"
	ProductIVCoinPackage ProductType = "ivCoinPackage"
)

func (p ProductType) String() string {
	return string(p)
}

func (p ProductType) IsValid() bool {
	switch p {
	case ProductInvitation, ProductIVCoinPackage:
		return true
	}
	return false
}

type TransactionStatusType string

const (
	TransactionStatusPending   TransactionStatusType = "pending"
	TransactionStatusCompleted TransactionStatusType = "completed"
)

func (t TransactionStatusType) String() string {
	return string(t)
}

func (t TransactionStatusType) IsValid() bool {
	switch t {
	case TransactionStatusPending, TransactionStatusCompleted:
		return true
	}
	return false
}

type PaymentMethodType string

const (
	PaymentMethodTransfer PaymentMethodType = "transfer"
	PaymentMethodIVCoin   PaymentMethodType = "ivCoin"
)

func (p PaymentMethodType) String() string {
	return string(p)
}

func (p PaymentMethodType) IsValid() bool {
	switch p {
	case PaymentMethodTransfer, PaymentMethodIVCoin:
		return true
	}
	return false
}

type Transaction struct {
	ID              uint                  `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductType     ProductType           `gorm:"type:varchar(50);not null" json:"product_type"`
	Status          TransactionStatusType `gorm:"type:varchar(50);not null;default:'pending'" json:"status"`
	PaymentMethod   PaymentMethodType     `gorm:"type:varchar(50);not null" json:"payment_method"`
	ReferenceNumber string                `gorm:"size:100;uniqueIndex;not null" json:"reference_number"`
	IDRPrice        uint                  `gorm:"not null" json:"idr_price"`
	IDRDiscount     uint                  `gorm:"not null" json:"idr_discount"`
	IDRTotalPrice   uint                  `gorm:"not null" json:"idr_total_price"`
	IVCPrice        uint                  `gorm:"not null" json:"ivc_price"`
	IVCDiscount     uint                  `gorm:"not null" json:"ivc_discount"`
	IVCTotalPrice   uint                  `gorm:"not null" json:"ivc_total_price"`
	ProductID       uint                  `gorm:"not null;index" json:"product_id"`
	Invitation      *Invitation           `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"invitation,omitempty"`
	IVCoinPackage   *IVCoinPackage        `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"iv_coin_package,omitempty"`

	UserID string `gorm:"size:36;not null;index" json:"user_id"`
	User   *User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
