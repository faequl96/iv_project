package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type User struct {
	ID        string                `gorm:"primaryKey;size:36" json:"id"` // Firebase UID
	UserName  string                `gorm:"size:100;uniqueIndex;not null" json:"user_name"`
	Email     string                `gorm:"size:255;uniqueIndex;not null" json:"email"`
	FullName  string                `gorm:"size:150;not null" json:"full_name"`
	IVCoin    *IVCoin               `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"iv_coin,omitempty"`
	CreatedAt time.Time             `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time             `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index" json:"-"`
}

type IVCoin struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    string    `gorm:"size:36;not null;index" json:"user_id"`
	User      *User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Balance   uint      `gorm:"not null;default:0" json:"balance"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type InvitationTheme struct {
	ID          uint                  `gorm:"primaryKey" json:"id"`
	Title       string                `gorm:"size:255;not null" json:"title"`
	NormalPrice int                   `gorm:"not null" json:"normal_price"`
	DiskonPrice int                   `gorm:"not null" json:"diskon_price"`
	Category    string                `gorm:"size:100;not null" json:"category"`
	Reviews     []*Review             `gorm:"foreignKey:InvitationThemeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"reviews,omitempty"`
	CreatedAt   time.Time             `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time             `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   soft_delete.DeletedAt `gorm:"index" json:"-"`
}

type Review struct {
	ID                uint             `gorm:"primaryKey" json:"id"`
	UserID            string           `gorm:"size:36;not null;index" json:"user_id"`
	User              *User            `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	InvitationThemeID uint             `gorm:"not null;index" json:"invitation_theme_id"`
	InvitationTheme   *InvitationTheme `gorm:"foreignKey:InvitationThemeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Star              int              `gorm:"not null" json:"star"`
	Comment           string           `gorm:"type:text;not null" json:"comment"`
	CreatedAt         time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
}

type Invitation struct {
	ID                uint                  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            string                `gorm:"size:36;not null;index" json:"user_id"`
	User              *User                 `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Status            string                `gorm:"size:10;not null;default:'pending'" json:"status"` // Gunakan validasi di aplikasi
	InvitationThemeID uint                  `gorm:"not null;index" json:"invitation_theme_id"`
	InvitationTheme   *InvitationTheme      `gorm:"foreignKey:InvitationThemeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"invitation_theme,omitempty"`
	InvitationData    *InvitationData       `gorm:"foreignKey:InvitationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"invitation_data,omitempty"`
	CreatedAt         time.Time             `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time             `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         soft_delete.DeletedAt `gorm:"index" json:"-"`
}

type InvitationData struct {
	ID           uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	InvitationID uint        `gorm:"not null;index" json:"invitation_id"`
	Invitation   *Invitation `gorm:"foreignKey:InvitationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	EventName    string      `gorm:"size:255;not null" json:"event_name"`
	EventDate    time.Time   `gorm:"not null" json:"event_date"`
	Location     string      `gorm:"size:255;not null" json:"location"`
	Gallery      []*Gallery  `gorm:"foreignKey:InvitationDataID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"gallery,omitempty"`
	CreatedAt    time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

type Gallery struct {
	ID               uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	InvitationDataID uint            `gorm:"not null;index" json:"invitation_data_id"`
	InvitationData   *InvitationData `gorm:"foreignKey:InvitationDataID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	ImageURL         string          `gorm:"size:255;not null" json:"image_url"`
}
