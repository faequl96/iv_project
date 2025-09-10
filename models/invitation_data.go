package models

import "time"

type InvitationData struct {
	ID             uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	CoverImageURL  string        `gorm:"size:255;not null" json:"main_image_url"`
	BrideID        uint          `gorm:"not null;index" json:"bride_id"`
	Bride          Bridegroom    `gorm:"foreignKey:BrideID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"bride"`
	GroomID        uint          `gorm:"not null;index" json:"groom_id"`
	Groom          Bridegroom    `gorm:"foreignKey:GroomID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"groom"`
	ContractID     uint          `gorm:"not null;index" json:"contract_id"`
	ContractEvent  Event         `gorm:"foreignKey:ContractID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"contract"`
	ReceptionID    uint          `gorm:"not null;index" json:"reception_id"`
	ReceptionEvent Event         `gorm:"foreignKey:ReceptionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"reception"`
	Gallery        *Gallery      `gorm:"foreignKey:InvitationDataID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"gallery,omitempty"`
	BankAccounts   []BankAccount `gorm:"foreignKey:InvitationDataID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"bank_accounts,omitempty"`

	InvitationID uint        `gorm:"not null;index" json:"invitation_id"`
	Invitation   *Invitation `gorm:"foreignKey:InvitationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
