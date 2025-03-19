package transaction_dto

import "iv_project/models"

type CreateTransactionRequest struct {
	ProductType   models.ProductType       `json:"product_type" binding:"required"`
	PaymentMethod models.PaymentMethodType `json:"payment_method" binding:"required"`
	ProductID     uint                     `json:"product_id" binding:"required"`
	UserID        string                   `json:"user_id" binding:"required"`
}

type UpdateTransactionRequest struct {
	Status models.TransactionStatusType `json:"status" binding:"omitempty"`
}
