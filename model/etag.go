package model

import (
	"time"
)

// Etag 是 EhTagTranslation/Database 的标签条目快照，
// 作为 decided_tags 手动匹配时的规范 tag 来源。
type Etag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	NS        string    `gorm:"type:varchar(16);uniqueIndex:uk_etag" json:"ns"`
	Key       string    `gorm:"type:varchar(128);uniqueIndex:uk_etag" json:"key"`
	Name      string    `gorm:"type:varchar(255);" json:"name"`
	Intro     string    `gorm:"type:text;" json:"intro"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli" json:"updated_at"`
}

func (Etag) TableName() string {
	return "etags"
}
