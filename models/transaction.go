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
	TransactionStatusConfirmed TransactionStatusType = "confirmed"
	TransactionStatusCanceled  TransactionStatusType = "canceled"
)

func (u TransactionStatusType) String() string {
	maps := map[TransactionStatusType]string{
		TransactionStatusPending:   "pending",
		TransactionStatusConfirmed: "confirmed",
		TransactionStatusCanceled:  "canceled",
	}
	return maps[u]
}

func StringToTransactionStatusType(value string) TransactionStatusType {
	maps := map[string]TransactionStatusType{
		"pending":   TransactionStatusPending,
		"confirmed": TransactionStatusConfirmed,
		"canceled":  TransactionStatusCanceled,
	}
	return maps[value]
}

type PaymentMethodType string

const (
	PaymentMethodIVCoin         PaymentMethodType = "iv_coin"
	PaymentMethodManualTransfer PaymentMethodType = "manual_transfer"
	PaymentMethodAutoTransfer   PaymentMethodType = "auto_transfer"
	PaymentMethodGopay          PaymentMethodType = "gopay"
)

func (u PaymentMethodType) String() string {
	maps := map[PaymentMethodType]string{
		PaymentMethodIVCoin:         "iv_coin",
		PaymentMethodManualTransfer: "manual_transfer",
		PaymentMethodAutoTransfer:   "auto_transfer",
		PaymentMethodGopay:          "gopay",
	}
	return maps[u]
}

func StringToPaymentMethodType(value string) PaymentMethodType {
	maps := map[string]PaymentMethodType{
		"iv_coin":         PaymentMethodIVCoin,
		"manual_transfer": PaymentMethodManualTransfer,
		"auto_transfer":   PaymentMethodAutoTransfer,
		"gopay":           PaymentMethodGopay,
	}
	return maps[value]
}

type Transaction struct {
	ID                     uint                  `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductType            ProductType           `gorm:"type:varchar(50);not null" json:"product_type"`
	ProductID              uint                  `gorm:"not null;index" json:"product_id"`
	ProductName            string                `gorm:"size:150;not null;index" json:"product_name"`
	Status                 TransactionStatusType `gorm:"type:varchar(50);not null;default:'pending'" json:"status"`
	PaymentMethod          PaymentMethodType     `gorm:"type:varchar(50);not null" json:"payment_method"`
	ReferenceNumber        string                `gorm:"size:100;uniqueIndex;not null" json:"reference_number"`
	IDRPrice               uint                  `gorm:"not null" json:"idr_price"`
	IDRDiscount            uint                  `gorm:"not null" json:"idr_discount"`
	IDRTotalPrice          uint                  `gorm:"not null" json:"idr_total_price"`
	IVCPrice               uint                  `gorm:"not null" json:"ivc_price"`
	IVCDiscount            uint                  `gorm:"not null" json:"ivc_discount"`
	IVCTotalPrice          uint                  `gorm:"not null" json:"ivc_total_price"`
	PaymentProofImageUrl   string                `gorm:"type:varchar(255)" json:"payment_proof_image_url"`
	VoucherCodeID          uint                  `gorm:"not null" json:"voucher_code_id"`
	VoucherCodeName        string                `gorm:"size:20;not null" json:"voucher_code_name"`
	IDRVoucherCodeDiscount uint                  `gorm:"not null" json:"idr_voucher_code_discount"`
	IVCVoucherCodeDiscount uint                  `gorm:"not null" json:"ivc_voucher_code_discount"`

	UserID string `gorm:"size:36;not null;index" json:"user_id"`
	User   *User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
