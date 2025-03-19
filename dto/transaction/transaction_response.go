package transaction_dto

import (
	invitation_dto "iv_project/dto/invitation"
	iv_coin_package_dto "iv_project/dto/iv_coin_package"
	"iv_project/models"
)

type TransactionResponse struct {
	ID              uint                                       `json:"id"`
	ProductType     models.ProductType                         `json:"product_type"`
	Status          models.TransactionStatusType               `json:"status"`
	PaymentMethod   models.PaymentMethodType                   `json:"payment_method"`
	ReferenceNumber string                                     `json:"reference_number"`
	IDRPrice        uint                                       `json:"idr_price"`
	IDRDiscount     uint                                       `json:"idr_discount"`
	IDRTotalPrice   uint                                       `json:"idr_total_price"`
	IVCPrice        uint                                       `json:"ivc_price"`
	IVCDiscount     uint                                       `json:"ivc_discount"`
	IVCTotalPrice   uint                                       `json:"ivc_total_price"`
	Invitation      *invitation_dto.InvitationResponse         `json:"invitation"`
	IVCoinPackage   *iv_coin_package_dto.IVCoinPackageResponse `json:"iv_coin_package"`
	CreatedAt       string                                     `json:"created_at"`
}
