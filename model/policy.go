package model

import "time"

type Policy[T any] struct {
	ID        uint      `gorm:"primaryKey"`
	Policy    T         `gorm:"type:json"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
}

type PlatformPolicy struct {
}
