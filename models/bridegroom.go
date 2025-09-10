package models

type Bridegroom struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Nickname    string `gorm:"type:varchar(255)" json:"nickname"`
	FullName    string `gorm:"type:varchar(255)" json:"full_name"`
	Title       string `gorm:"type:varchar(255)" json:"title"`
	ImageURL    string `gorm:"type:varchar(255)" json:"image_url"`
	FatherName  string `gorm:"type:varchar(255)" json:"father_name"`
	FatherTitle string `gorm:"type:varchar(255)" json:"father_title"`
	MotherName  string `gorm:"type:varchar(255)" json:"mother_name"`
	MotherTitle string `gorm:"type:varchar(255)" json:"mother_title"`
}
