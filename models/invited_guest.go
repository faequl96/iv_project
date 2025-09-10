package models

type InvitedGuest struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	InvitationID uint   `gorm:"not null;index" json:"invitation_id"`
	NameInstance string `gorm:"size:255;not null" json:"name_instance"`
	Name         string `gorm:"size:255;not null" json:"name"`
	Instance     string `gorm:"size:255;not null" json:"instance"`
	Nickname     string `gorm:"size:255;not null" json:"nickname"`
	Avatar       string `gorm:"size:255;not null" json:"avatar"`
	Attendance   string `gorm:"size:255;not null" json:"attendance"`
}
