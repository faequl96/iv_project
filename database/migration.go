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
		&models.InvitationTheme{},
		&models.Invitation{},
		&models.InvitationData{},
		&models.Gallery{},
		&models.Review{},
	)

	if err != nil {
		fmt.Println(err)
		panic("Migration Failed")
	}

	fmt.Println("Migration Success")
}
