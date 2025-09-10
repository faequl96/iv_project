package rsvp_dto

type RSVPResponse struct {
	ID             uint   `json:"id"`
	InvitationID   uint   `json:"invitation_id"`
	InvitedGuestID uint   `json:"invited_guest_id"`
	Nickname       string `json:"nickname"`
	Avatar         string `json:"avatar"`
	Invited        bool   `json:"invited"`
	Attendance     string `json:"attendance"`
	Message        string `json:"message"`
}
