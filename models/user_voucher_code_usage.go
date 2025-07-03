package models

type UserVoucherCodeUsage struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	UserID        string `json:"user_id"`
	VoucherCodeID uint   `json:"voucher_code_id"`
	UsageCount    int    `json:"usage_count"`

	User    User        `gorm:"foreignKey:UserID" json:"-"`
	Voucher VoucherCode `gorm:"foreignKey:VoucherCodeID" json:"-"`
}
