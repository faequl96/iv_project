package models

import "time"

type IVCoin struct {
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Balance           uint      `gorm:"not null;default:0" json:"balance"`
	AdMobMarker       uint      `gorm:"not null;default:0" json:"ad_mob_marker"`
	AdMobLastUpdateAt time.Time `gorm:"column:ad_mob_last_update_at" json:"ad_mob_last_update_at"`

	UserID string `gorm:"uniqueIndex;not null" json:"user_id"`
	User   *User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}
