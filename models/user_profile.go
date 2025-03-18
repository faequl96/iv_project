package models

type UserProfile struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string `gorm:"size:120;not null" json:"first_name"`
	LastName  string `gorm:"size:120;not null" json:"last_name"`

	UserID string `gorm:"size:36;not null;index" json:"user_id"`
	User   *User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}
