package invitation_data_dto

import (
	bank_account_dto "iv_project/dto/bank_account"
	bridegroom_dto "iv_project/dto/bridegroom"
	event_dto "iv_project/dto/event"
)

type CreateInvitationDataRequest struct {
	Bride          bridegroom_dto.UpdateBridegroomRequest `json:"bride" validate:"required"`
	Groom          bridegroom_dto.UpdateBridegroomRequest `json:"groom" validate:"required"`
	ContractEvent  event_dto.UpdateEventRequest           `json:"contract_event" validate:"required"`
	ReceptionEvent event_dto.UpdateEventRequest           `json:"reception_event" validate:"required"`
	BankAccounts   []bank_account_dto.BankAccountRequest  `json:"bank_accounts" validate:"required"`
}

type UpdateInvitationDataRequest struct {
	Bride          bridegroom_dto.UpdateBridegroomRequest `json:"bride"`
	Groom          bridegroom_dto.UpdateBridegroomRequest `json:"groom"`
	ContractEvent  event_dto.UpdateEventRequest           `json:"contract_event"`
	ReceptionEvent event_dto.UpdateEventRequest           `json:"reception_event"`
	BankAccounts   []bank_account_dto.BankAccountRequest  `json:"bank_accounts"`
}
