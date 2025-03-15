package models

import (
	"time"
)

type IVCoin struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    string    `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user_id"`
	Balance   int       `gorm:"not null" json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
