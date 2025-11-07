package model

import "time"

type Brand struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	VNDBID    string    `gorm:"type:varchar(255);index:brand_vndbid_idx" json:"vndb_id"`
	Name      string    `gorm:"type:varchar(255);unique" json:"name"`
	Links     Links     `gorm:"type:json" json:"link"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
}

func (Brand) TableName() string {
	return "brands"
}
