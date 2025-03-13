package models

import (
	"time"
)

type IVCoint struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    string    `gorm:"type:varchar(255);uniqueIndex;constraint:OnDelete:CASCADE;not null" json:"user_id"`
	Balance   int       `gorm:"not null" json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
