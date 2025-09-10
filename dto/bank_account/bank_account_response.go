package bank_account_dto

type BankAccountResponse struct {
	ID          uint   `json:"id"`
	BankName    string `json:"bank_name"`
	AccountName string `json:"account_name"`
	Number      string `json:"number"`
}
