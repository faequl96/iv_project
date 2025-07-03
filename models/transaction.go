package models

import (
	"time"
)

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
	TransactionStatusCreated   TransactionStatusType = "created"
	TransactionStatusPending   TransactionStatusType = "pending"
	TransactionStatusConfirmed TransactionStatusType = "confirmed"
	TransactionStatusCanceled  TransactionStatusType = "canceled"
)

func (u TransactionStatusType) String() string {
	maps := map[TransactionStatusType]string{
		TransactionStatusCreated:   "created",
		TransactionStatusPending:   "pending",
		TransactionStatusConfirmed: "confirmed",
		TransactionStatusCanceled:  "canceled",
	}
	return maps[u]
}

func StringToTransactionStatusType(value string) TransactionStatusType {
	maps := map[string]TransactionStatusType{
		"created":   TransactionStatusCreated,
		"pending":   TransactionStatusPending,
		"confirmed": TransactionStatusConfirmed,
		"canceled":  TransactionStatusCanceled,
	}
	return maps[value]
}

type MidtransTransactionStatusType string

const (
	MidtransTransactionStatusPending    MidtransTransactionStatusType = "pending"
	MidtransTransactionStatusUnknown    MidtransTransactionStatusType = "unknown"
	MidtransTransactionStatusSettlement MidtransTransactionStatusType = "settlement"
	MidtransTransactionStatusCapture    MidtransTransactionStatusType = "capture"
	MidtransTransactionStatusExpire     MidtransTransactionStatusType = "expire"
	MidtransTransactionStatusCancel     MidtransTransactionStatusType = "cancel"
	MidtransTransactionStatusDeny       MidtransTransactionStatusType = "deny"
)

func (u MidtransTransactionStatusType) String() string {
	maps := map[MidtransTransactionStatusType]string{
		MidtransTransactionStatusUnknown:    "unknown",
		MidtransTransactionStatusPending:    "pending",
		MidtransTransactionStatusSettlement: "settlement",
		MidtransTransactionStatusCapture:    "capture",
		MidtransTransactionStatusExpire:     "expire",
		MidtransTransactionStatusCancel:     "cancel",
		MidtransTransactionStatusDeny:       "deny",
	}
	return maps[u]
}

func StringToMidtransTransactionStatusType(value string) MidtransTransactionStatusType {
	maps := map[string]MidtransTransactionStatusType{
		"unknown":    MidtransTransactionStatusUnknown,
		"pending":    MidtransTransactionStatusPending,
		"settlement": MidtransTransactionStatusSettlement,
		"capture":    MidtransTransactionStatusCapture,
		"expire":     MidtransTransactionStatusExpire,
		"cancel":     MidtransTransactionStatusCancel,
		"deny":       MidtransTransactionStatusDeny,
	}
	return maps[value]
}

type PaymentMethodType string

const (
	PaymentMethodIVCoin       PaymentMethodType = "iv_coin"
	PaymentMethodGopay        PaymentMethodType = "gopay"
	PaymentMethodQRIS         PaymentMethodType = "qris"
	PaymentMethodBankTransfer PaymentMethodType = "bank_transfer"
)

func (u PaymentMethodType) String() string {
	maps := map[PaymentMethodType]string{
		PaymentMethodIVCoin:       "iv_coin",
		PaymentMethodGopay:        "gopay",
		PaymentMethodQRIS:         "qris",
		PaymentMethodBankTransfer: "bank_transfer",
	}
	return maps[u]
}

func StringToPaymentMethodType(value string) PaymentMethodType {
	maps := map[string]PaymentMethodType{
		"iv_coin":       PaymentMethodIVCoin,
		"gopay":         PaymentMethodGopay,
		"qris":          PaymentMethodQRIS,
		"bank_transfer": PaymentMethodBankTransfer,
	}
	return maps[value]
}

type Transaction struct {
	ID                     string                        `gorm:"primaryKey;size:36" json:"id"`
	TransactionCode        string                        `gorm:"size:20;uniqueIndex;not null" json:"transaction_code"`
	ProductType            ProductType                   `gorm:"type:varchar(50);not null" json:"product_type"`
	ProductID              uint                          `gorm:"not null;index" json:"product_id"`
	ProductName            string                        `gorm:"size:150;not null;index" json:"product_name"`
	ProductDescription     string                        `gorm:"not null;index" json:"product_description"`
	Status                 TransactionStatusType         `gorm:"type:varchar(50);not null;default:'create'" json:"status"`
	PaymentMethod          PaymentMethodType             `gorm:"type:varchar(50);not null" json:"payment_method"`
	ReferenceNumber        string                        `gorm:"size:100;uniqueIndex;not null" json:"reference_number"`
	PaymentURL             string                        `gorm:"not null" json:"payment_url"`
	Acquirer               string                        `gorm:"not null" json:"acquirer"`
	MidtransStatus         MidtransTransactionStatusType `gorm:"type:varchar(50);not null;default:'unknown'" json:"midtrans_status"`
	TimeLimitAt            *time.Time                    `gorm:"column:time_limit_at" json:"time_limit_at"`
	IDRPrice               uint                          `gorm:"not null" json:"idr_price"`
	IDRDiscount            uint                          `gorm:"not null" json:"idr_discount"`
	IDRTotalPrice          uint                          `gorm:"not null" json:"idr_total_price"`
	IVCPrice               uint                          `gorm:"not null" json:"ivc_price"`
	IVCDiscount            uint                          `gorm:"not null" json:"ivc_discount"`
	IVCTotalPrice          uint                          `gorm:"not null" json:"ivc_total_price"`
	VoucherCodeID          uint                          `gorm:"not null" json:"voucher_code_id"`
	VoucherCodeName        string                        `gorm:"size:20;not null" json:"voucher_code_name"`
	IDRVoucherCodeDiscount uint                          `gorm:"not null" json:"idr_voucher_code_discount"`
	IVCVoucherCodeDiscount uint                          `gorm:"not null" json:"ivc_voucher_code_discount"`
	PaymentProofImageUrl   string                        `gorm:"type:varchar(255)" json:"payment_proof_image_url"`

	UserID string `gorm:"size:36;not null;index" json:"user_id"`
	User   *User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
