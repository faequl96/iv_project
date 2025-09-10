package invitation_data_dto

import (
	bank_account_dto "iv_project/dto/bank_account"
	bridegroom_dto "iv_project/dto/bridegroom"
	event_dto "iv_project/dto/event"
	gallery_dto "iv_project/dto/gallery"
)

type InvitationDataResponse struct {
	ID             uint                                   `json:"id"`
	CoverImageURL  string                                 `json:"cover_image_url"`
	Bride          bridegroom_dto.BridegroomResponse      `json:"bride"`
	Groom          bridegroom_dto.BridegroomResponse      `json:"groom"`
	ContractEvent  event_dto.EventResponse                `json:"contract_event"`
	ReceptionEvent event_dto.EventResponse                `json:"reception_event"`
	Gallery        *gallery_dto.GalleryResponse           `json:"gallery"`
	BankAccounts   []bank_account_dto.BankAccountResponse `json:"bank_accounts"`
}
