package transaction_confirmation_dto

import "iv_project/models"

type TransactionConfirmationRequest struct {
	Status models.TransactionStatusType `json:"status" validate:"required"`
}
