package database

import (
	"fmt"
	"iv_project/models"
	"iv_project/pkg/mysql"
)

func RunMigration() {
	err := mysql.DB.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.IVCoin{},
		&models.Category{},
		&models.DiscountCategory{},
		&models.IVCoinPackage{},
		&models.InvitationTheme{},
		&models.Review{},
		&models.Invitation{},
		&models.InvitedGuest{},
		&models.RSVP{},
		&models.InvitationData{},
		&models.Bridegroom{},
		&models.Event{},
		&models.BankAccount{},
		&models.Gallery{},
		&models.Transaction{},
		&models.VoucherCode{},
		&models.UserVoucherCodeUsage{},
	)

	if err != nil {
		fmt.Println(err)
		panic("Migration Failed")
	}

	fmt.Println("Migration Success")
}
