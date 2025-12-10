package model

import "time"

type Series struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);unique" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`

	Games int64 `gorm:"->"`
}

func (Series) TableName() string {
	return "series"
}
