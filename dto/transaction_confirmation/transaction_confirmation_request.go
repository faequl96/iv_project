package transaction_confirmation_dto

import "iv_project/models"

type UpdateTransactionConfirmationRequest struct {
	Status models.TransactionStatusType `json:"status"`
}
