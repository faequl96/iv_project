package models

import (
	"time"
)

type User struct {
	ID        string    `gorm:"primaryKey;size:36" json:"id"`
	UserName  string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"user_name"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	FullName  string    `gorm:"type:varchar(255);not null" json:"full_name"`
	IVCoint   IVCoint   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"iv_coint"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
