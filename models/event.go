package models

import "time"

type Event struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	StartTime time.Time  `gorm:"column:start_time" json:"start_time"`
	EndTime   *time.Time `gorm:"column:end_time" json:"end_time"`
	Place     string     `gorm:"type:varchar(255)" json:"place"`
	Address   string     `gorm:"type:varchar(255)" json:"address"`
	MapsURL   string     `gorm:"type:varchar(255)" json:"maps_url"`
}
