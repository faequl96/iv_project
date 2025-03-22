package models

import "time"

type ProductType string

const (
	ProductInvitation    ProductType = "invitation"
	ProductIVCoinPackage ProductType = "iv_coin_package"
)

func (u ProductType) String() string {
	maps := map[ProductType]string{
		ProductInvitation:    "invitation",
		ProductIVCoinPackage: "iv_coin_package",
	}
	return maps[u]
}

func StringToProductType(value string) ProductType {
	maps := map[string]ProductType{
		"invitation":      ProductInvitation,
		"iv_coin_package": ProductIVCoinPackage,
	}
	return maps[value]
}

type TransactionStatusType string

const (
	TransactionStatusPending   TransactionStatusType = "pending"
	TransactionStatusCompleted TransactionStatusType = "completed"
)

func (u TransactionStatusType) String() string {
	maps := map[TransactionStatusType]string{
		TransactionStatusPending:   "pending",
		TransactionStatusCompleted: "completed",
	}
	return maps[u]
}

func StringToTransactionStatusType(value string) TransactionStatusType {
	maps := map[string]TransactionStatusType{
		"pending":   TransactionStatusPending,
		"completed": TransactionStatusCompleted,
	}
	return maps[value]
}

type PaymentMethodType string

const (
	PaymentMethodTransfer PaymentMethodType = "transfer"
	PaymentMethodIVCoin   PaymentMethodType = "iv_coin"
)

func (u PaymentMethodType) String() string {
	maps := map[PaymentMethodType]string{
		PaymentMethodTransfer: "transfer",
		PaymentMethodIVCoin:   "iv_coin",
	}
	return maps[u]
}

func StringToPaymentMethodType(value string) PaymentMethodType {
	maps := map[string]PaymentMethodType{
		"transfer": PaymentMethodTransfer,
		"iv_coin":  PaymentMethodIVCoin,
	}
	return maps[value]
}

type Transaction struct {
	ID              uint                  `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductType     ProductType           `gorm:"type:varchar(50);not null" json:"product_type"`
	ProductName     string                `gorm:"size:150;not null;index" json:"product_name"`
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
