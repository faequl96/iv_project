package transaction_dto

import "iv_project/models"

type CreateTransactionRequest struct {
	ProductType   models.ProductType       `json:"product_type" validate:"required"`
	PaymentMethod models.PaymentMethodType `json:"payment_method" validate:"required"`

	ProductID uint   `json:"product_id" validate:"required"`
	UserID    string `json:"user_id" validate:"required"`
}

type UpdateTransactionRequest struct {
	Status models.TransactionStatusType `json:"status"`
}
