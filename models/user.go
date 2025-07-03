package models

import "time"

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

func StringToUserRoleType(value string) UserRoleType {
	maps := map[string]UserRoleType{
		"super_admin": UserRoleSuperAdmin,
		"admin":       UserRoleAdmin,
		"user":        UserRoleUser,
	}
	return maps[value]
}

type User struct {
	ID           string        `gorm:"primaryKey;size:36" json:"id"` // Firebase UID
	UnixID       string        `gorm:"primaryKey;size:100;not null" json:"unix_id"`
	Role         UserRoleType  `gorm:"type:varchar(50);not null;default:'user'" json:"role"`
	UserProfile  *UserProfile  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user_profile,omitempty"`
	IVCoin       *IVCoin       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"iv_coin,omitempty"`
	VoucherCodes []VoucherCode `gorm:"many2many:voucher_code_users;" json:"-"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
