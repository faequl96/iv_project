package invited_guest_dto

type CreateInvitedGuestRequest struct {
	InvitationID uint   `json:"invitation_id" validate:"required"`
	NameInstance string `json:"name_instance" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Instance     string `json:"instance" validate:"required"`
	Nickname     string `json:"nickname" validate:"required"`
	Avatar       string `json:"avatar" validate:"required"`
	Attendance   string `json:"attendance" validate:"required"`
}

type UpdateInvitedGuestRequest struct {
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	Attendance string `json:"attendance"`
}
