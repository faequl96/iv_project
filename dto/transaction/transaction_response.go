package transaction_dto

import "iv_project/models"

type TransactionResponse struct {
	ID              uint                         `json:"id"`
	ProductType     models.ProductType           `json:"product_type"`
	ProductName     string                       `json:"product_name"`
	Status          models.TransactionStatusType `json:"status"`
	PaymentMethod   models.PaymentMethodType     `json:"payment_method"`
	ReferenceNumber string                       `json:"reference_number"`
	IDRPrice        uint                         `json:"idr_price"`
	IDRDiscount     uint                         `json:"idr_discount"`
	IDRTotalPrice   uint                         `json:"idr_total_price"`
	IVCPrice        uint                         `json:"ivc_price"`
	IVCDiscount     uint                         `json:"ivc_discount"`
	IVCTotalPrice   uint                         `json:"ivc_total_price"`
	CreatedAt       string                       `json:"created_at"`
}
