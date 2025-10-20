package model

import "time"

type Tag struct {
	ID        uint      `gorm:"primaryKey;type:bigint" json:"id"`
	Name      string    `gorm:"varchar(255);unique" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
}

func (Tag) TableName() string {
	return "tags"
}
