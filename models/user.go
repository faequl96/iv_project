package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type UserRoleType string

const (
	UserRoleSuperAdmin UserRoleType = "super_admin"
	UserRoleAdmin      UserRoleType = "admin"
	UserRoleUser       UserRoleType = "user"
)

func (u UserRoleType) String() string {
	maps := map[UserRoleType]string{
		UserRoleSuperAdmin: "super_admin",
		UserRoleAdmin:      "admin",
		UserRoleUser:       "user",
	}
	return maps[u]
}

func StringToUserRole(value string) UserRoleType {
	maps := map[string]UserRoleType{
		"super_admin": UserRoleSuperAdmin,
		"admin":       UserRoleAdmin,
		"user":        UserRoleUser,
	}
	return maps[value]
}

type User struct {
	ID          string       `gorm:"primaryKey;size:36" json:"id"` // Firebase UID
	Email       string       `gorm:"size:100" json:"email"`
	Role        UserRoleType `gorm:"type:varchar(50);not null;default:'user'" json:"role"`
	UserProfile *UserProfile `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user_profile,omitempty"`
	IVCoin      *IVCoin      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"iv_coin,omitempty"`

	CreatedAt time.Time             `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time             `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index" json:"-"`
}
