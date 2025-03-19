package invitation_dto

import (
	invitation_data_dto "iv_project/dto/invitation_data"
	"iv_project/models"
)

type CreateInvitationRequest struct {
	UserID            string                                          `json:"user_id" validate:"required"`
	InvitationThemeID uint                                            `json:"invitation_theme_id" validate:"required"`
	InvitationData    invitation_data_dto.CreateInvitationDataRequest `json:"invitation_data"`
}

type UpdateInvitationRequest struct {
	Status         models.InvitationStatusType                     `json:"status"`
	InvitationData invitation_data_dto.UpdateInvitationDataRequest `json:"invitation_data"`
}
