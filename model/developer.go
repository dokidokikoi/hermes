package model

import "time"

type Developer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `gorm:"type:varchar(255);unique" json:"uuid"`
	Name      string    `gorm:"type:varchar(255)" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
}

func (Developer) TableName() string {
	return "developers"
}
