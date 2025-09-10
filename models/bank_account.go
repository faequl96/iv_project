package models

type BankAccount struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	BankName    string `gorm:"type:varchar(255)" json:"bank_name"`
	AccountName string `gorm:"type:varchar(255)" json:"account_name"`
	Number      string `gorm:"type:varchar(255)" json:"number"`

	InvitationDataID uint            `gorm:"not null;index" json:"invitation_data_id"`
	InvitationData   *InvitationData `gorm:"foreignKey:InvitationDataID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}
