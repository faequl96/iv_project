package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AuthID    string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"auth_id"`
	UserName  string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"user_name"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	FullName  string    `gorm:"type:varchar(255);not null" json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
