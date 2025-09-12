package model

import (
	"time"

	"github.com/lib/pq"
)

type Character struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UUID      string         `gorm:"type:varchar(255);unique" json:"uuid"`
	Name      string         `gorm:"type:varchar(255)" json:"name"`
	Alias     pq.StringArray `gorm:"type:json" json:"alias"`
	Gender    Gender         `gorm:"type:varchar(10)" json:"gender"`
	Summary   string         `json:"summary"`
	Cover     string         `gorm:"type:varchar(512)" json:"cover"`
	Images    pq.StringArray `gorm:"type:json" json:"images"`
	CV        Person         `gorm:"foreignKey:PersonID" json:"cv"`
	PersonID  uint           `gorm:"default:null" json:"-"`
	Tags      []Tag          `gorm:"many2many:character_tag;" json:"tags"`
	CreatedAt time.Time      `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:milli"`
}

func (Character) TableName() string {
	return "characters"
}

type CharacterTag struct {
	CharacterID uint `gorm:"primaryKey"`
	TagID       uint `gorm:"primaryKey"`
}

func (CharacterTag) TableName() string {
	return "character_tag"
}
