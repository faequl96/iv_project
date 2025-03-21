package models

type UserProfile struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string `gorm:"size:120;not null" json:"first_name"`
	LastName  string `gorm:"size:120;not null" json:"last_name"`

	UserID string `gorm:"uniqueIndex;not null" json:"user_id"`
}
