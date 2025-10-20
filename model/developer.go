package model

import "time"

type Brand struct {
	ID        uint      `gorm:"primaryKey;type:bigint" json:"id"`
	Name      string    `gorm:"type:varchar(255);unique" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
}

func (Brand) TableName() string {
	return "brands"
}
