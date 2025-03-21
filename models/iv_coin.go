package models

type IVCoin struct {
	ID      uint `gorm:"primaryKey;autoIncrement" json:"id"`
	Balance uint `gorm:"not null;default:0" json:"balance"`

	UserID string `gorm:"uniqueIndex;not null" json:"user_id"`
	User   *User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}
