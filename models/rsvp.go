package models

type RSVP struct {
	ID             uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	InvitationID   uint   `gorm:"not null;index" json:"invitation_id"`
	InvitedGuestID uint   `gorm:"not null;index" json:"invited_guest_id"`
	Nickname       string `gorm:"size:255;not null" json:"nickname"`
	Avatar         string `gorm:"size:255;not null" json:"avatar"`
	Invited        bool   `json:"invited"`
	Attendance     string `gorm:"size:255;not null" json:"attendance"`
	Message        string `gorm:"size:255;not null" json:"message"`
}
