package transaction_dto

import (
	"iv_project/models"
	"time"
)

type TransactionResponse struct {
	ID                     string                               `json:"id"`
	TransactionCode        string                               `json:"transaction_code"`
	ProductType            models.ProductType                   `json:"product_type"`
	ProductName            string                               `json:"product_name"`
	ProductDescription     string                               `json:"product_description"`
	Status                 models.TransactionStatusType         `json:"status"`
	PaymentMethod          models.PaymentMethodType             `json:"payment_method"`
	PaymentURL             string                               `json:"payment_url"`
	ReferenceNumber        string                               `json:"reference_number"`
	Acquirer               string                               `json:"acquirer"`
	MidtransStatus         models.MidtransTransactionStatusType `json:"midtrans_status"`
	TimeLimitAt            *time.Time                           `json:"time_limit_at"`
	IDRPrice               uint                                 `json:"idr_price"`
	IDRDiscount            uint                                 `json:"idr_discount"`
	IDRTotalPrice          uint                                 `json:"idr_total_price"`
	IVCPrice               uint                                 `json:"ivc_price"`
	IVCDiscount            uint                                 `json:"ivc_discount"`
	IVCTotalPrice          uint                                 `json:"ivc_total_price"`
	VoucherCodeID          uint                                 `json:"voucher_code_id"`
	VoucherCodeName        string                               `json:"voucher_code_name"`
	IDRVoucherCodeDiscount uint                                 `json:"idr_voucher_code_discount"`
	IVCVoucherCodeDiscount uint                                 `json:"ivc_voucher_code_discount"`
	PaymentProofImageUrl   string                               `json:"payment_proof_image_url"`
	CreatedAt              string                               `json:"created_at"`
}
