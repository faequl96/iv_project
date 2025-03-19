package models

type IVCoin struct {
	ID      uint `gorm:"primaryKey;autoIncrement" json:"id"`
	Balance uint `gorm:"not null;default:0" json:"balance"`

	UserID string `gorm:"size:36;not null;index" json:"user_id"`
}
