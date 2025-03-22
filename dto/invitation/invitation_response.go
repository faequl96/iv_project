package invitation_dto

import (
	invitation_data_dto "iv_project/dto/invitation_data"
	invitation_theme_dto "iv_project/dto/invitation_theme"
	"iv_project/models"
)

type InvitationResponse struct {
	ID              uint                                          `json:"id"`
	Status          models.InvitationStatusType                   `json:"status"`
	InvitationTheme *invitation_theme_dto.InvitationThemeResponse `json:"invitation_theme"`
	InvitationData  *invitation_data_dto.InvitationDataResponse   `json:"invitation_data"`
}
