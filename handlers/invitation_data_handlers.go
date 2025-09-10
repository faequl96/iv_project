package handlers

import (
	bank_account_dto "iv_project/dto/bank_account"
	invitation_data_dto "iv_project/dto/invitation_data"
	"iv_project/models"
)

func ConvertToInvitationDataResponse(invitationData *models.InvitationData) invitation_data_dto.InvitationDataResponse {
	invitationDataResponse := invitation_data_dto.InvitationDataResponse{
		ID:             invitationData.ID,
		CoverImageURL:  invitationData.CoverImageURL,
		Bride:          ConvertToBridegroomResponse(invitationData.Bride),
		Groom:          ConvertToBridegroomResponse(invitationData.Groom),
		ContractEvent:  ConvertToEventResponse(invitationData.ContractEvent),
		ReceptionEvent: ConvertToEventResponse(invitationData.ReceptionEvent),
		BankAccounts:   []bank_account_dto.BankAccountResponse{},
	}

	var bankAccountResponses []bank_account_dto.BankAccountResponse
	for _, bankAccount := range invitationData.BankAccounts {
		bankAccountResponses = append(bankAccountResponses, ConvertToBankAccountResponse(bankAccount))
	}
	if len(bankAccountResponses) != 0 {
		invitationDataResponse.BankAccounts = bankAccountResponses
	}

	if invitationData.Gallery != nil {
		galleryResponse := ConvertToGalleryResponse(invitationData.Gallery)
		invitationDataResponse.Gallery = &galleryResponse
	}

	return invitationDataResponse
}
