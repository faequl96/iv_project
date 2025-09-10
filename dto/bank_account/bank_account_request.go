package bank_account_dto

type BankAccountRequest struct {
	BankName    string `json:"bank_name" validate:"required"`
	AccountName string `json:"account_name" validate:"required"`
	Number      string `json:"number" validate:"required"`
}
