package invitation_dto

import (
	invitation_data_dto "iv_project/dto/invitation_data"
	"iv_project/models"
)

type InvitationResponse struct {
	ID                  uint                                        `json:"id"`
	Status              models.InvitationStatusType                 `json:"status"`
	InvitationThemeID   uint                                        `json:"invitation_theme_id"`
	InvitationThemeName string                                      `json:"invitation_theme_name"`
	InvitationData      *invitation_data_dto.InvitationDataResponse `json:"invitation_data"`
}
