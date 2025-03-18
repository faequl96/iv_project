package invitation_dto

import (
	invitation_data_dto "iv_project/dto/invitation_data"
	invitation_theme_dto "iv_project/dto/invitation_theme"
	user_dto "iv_project/dto/user"
)

type InvitationResponse struct {
	ID              uint                                          `json:"id"`
	Status          string                                        `json:"status"`
	User            *user_dto.UserResponse                        `json:"user"`
	InvitationTheme *invitation_theme_dto.InvitationThemeResponse `json:"invitation_theme"`
	InvitationData  *invitation_data_dto.InvitationDataResponse   `json:"invitation_data"`
}
