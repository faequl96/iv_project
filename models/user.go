package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type User struct {
	ID        string                `gorm:"primaryKey;size:36" json:"id"` // Firebase UID
	UserName  string                `gorm:"type:varchar(100);uniqueIndex;not null" json:"user_name"`
	Email     string                `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	FullName  string                `gorm:"type:varchar(255);not null" json:"full_name"`
	IVCoin    *IVCoin               `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"iv_coin,omitempty"`
	CreatedAt time.Time             `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time             `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
