package handlers

import (
	bank_account_dto "iv_project/dto/bank_account"
	"iv_project/models"
)

func ConvertToBankAccountResponse(bankAccount models.BankAccount) bank_account_dto.BankAccountResponse {
	return bank_account_dto.BankAccountResponse{
		ID:          bankAccount.ID,
		BankName:    bankAccount.BankName,
		AccountName: bankAccount.AccountName,
		Number:      bankAccount.Number,
	}
}
