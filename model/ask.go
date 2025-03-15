package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"primaryKey;size:36" json:"id"` // Firebase UID
	UserName  string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"user_name"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	FullName  string         `gorm:"type:varchar(255);not null" json:"full_name"`
	IVCoin    *IVCoin        `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"iv_coin,omitempty"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type IVCoin struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    string    `gorm:"not null;index" json:"user_id"`
	User      *User     `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	Balance   uint      `gorm:"not null;default:0" json:"balance"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type InvitationTheme struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	NormalPrice int            `gorm:"not null" json:"normal_price"`
	DiskonPrice int            `gorm:"not null" json:"diskon_price"`
	Category    string         `gorm:"type:varchar(100);not null" json:"category"`
	Reviews     []Review       `gorm:"foreignKey:InvitationThemeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"reviews,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type Review struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	UserID            string          `gorm:"not null" json:"user_id"`
	User              *User           `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	InvitationThemeID uint            `gorm:"not null;index" json:"invitation_theme_id"`
	InvitationTheme   InvitationTheme `gorm:"foreignKey:InvitationThemeID;references:ID" json:"-"`
	Star              int             `gorm:"not null" json:"star"`
	Comment           string          `gorm:"type:text;not null" json:"comment"`
	CreatedAt         time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

type Invitation struct {
	ID                uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            string          `gorm:"not null;index" json:"user_id"`
	User              *User           `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Status            string          `gorm:"type:enum('pending','approved','rejected');not null;default:'pending'" json:"status"`
	InvitationThemeID uint            `gorm:"not null;index" json:"invitation_theme_id"`
	InvitationTheme   InvitationTheme `gorm:"foreignKey:InvitationThemeID;references:ID" json:"invitation_theme"`
	InvitationData    *InvitationData `gorm:"foreignKey:InvitationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"invitation_data,omitempty"`
	CreatedAt         time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
}

type InvitationData struct {
	ID           uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	InvitationID uint        `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"invitation_id"`
	Invitation   *Invitation `gorm:"foreignKey:InvitationID;references:ID" json:"-"`
	EventName    string      `gorm:"type:varchar(255);not null" json:"event_name"`
	EventDate    time.Time   `gorm:"not null" json:"event_date"`
	Location     string      `gorm:"type:varchar(255);not null" json:"location"`
	Gallery      []Gallery   `gorm:"foreignKey:InvitationDataID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"gallery,omitempty"`
	CreatedAt    time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

type Gallery struct {
	ID               uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	InvitationDataID uint            `gorm:"not null;index" json:"invitation_data_id"`
	InvitationData   *InvitationData `gorm:"foreignKey:InvitationDataID;references:ID" json:"-"`
	ImageURL         string          `gorm:"type:varchar(255);not null" json:"image_url"`
}
