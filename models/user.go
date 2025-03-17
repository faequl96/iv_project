package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type User struct {
	ID          string                `gorm:"primaryKey;size:36" json:"id"` // Firebase UID
	Email       string                `gorm:"size:255;uniqueIndex;not null" json:"email"`
	UserName    string                `gorm:"size:60;uniqueIndex;not null" json:"user_name"`
	UserProfile *UserProfile          `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user_profile,omitempty"`
	IVCoin      *IVCoin               `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"iv_coin,omitempty"`
	CreatedAt   time.Time             `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time             `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   soft_delete.DeletedAt `gorm:"index" json:"-"`
}
