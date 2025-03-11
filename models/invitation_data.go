package models

import "time"

type InvitationData struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID   uint      `gorm:"not null;uniqueIndex;constraint:OnDelete:CASCADE;" json:"order_id"`
	EventName string    `gorm:"type:varchar(255);not null" json:"event_name"`
	EventDate time.Time `gorm:"not null" json:"event_date"`
	Location  string    `gorm:"type:varchar(255);not null" json:"location"`
}
