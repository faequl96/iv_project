package invited_guest_dto

type InvitedGuestResponse struct {
	ID           uint   `json:"id"`
	InvitationID uint   `json:"invitation_id"`
	NameInstance string `json:"name_instance"`
	Name         string `json:"name"`
	Instance     string `json:"instance"`
	Nickname     string `json:"nickname"`
	Avatar       string `json:"avatar"`
	Attendance   string `json:"attendance"`
}
