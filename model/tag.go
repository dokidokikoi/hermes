package model

import "time"

type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"unique" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
}
