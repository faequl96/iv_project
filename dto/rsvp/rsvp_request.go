package rsvp_dto

type RSVPRequest struct {
	InvitationID   uint   `json:"invitation_id" validate:"required"`
	InvitedGuestID uint   `json:"invited_guest_id" validate:"required"`
	Nickname       string `json:"nickname" validate:"required"`
	Avatar         string `json:"avatar" validate:"required"`
	Invited        bool   `json:"invited" validate:"required"`
	Attendance     string `json:"attendance" validate:"required"`
	Message        string `json:"message" validate:"required"`
}
