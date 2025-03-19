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
		&models.IVCoinPackage{},
		&models.InvitationTheme{},
		&models.Category{},
		&models.Review{},
		&models.Invitation{},
		&models.InvitationData{},
		&models.Gallery{},
		&models.Transaction{},
	)

	if err != nil {
		fmt.Println(err)
		panic("Migration Failed")
	}

	fmt.Println("Migration Success")
}
