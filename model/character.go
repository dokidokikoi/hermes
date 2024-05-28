package model

import "time"

type Character struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Alias     []string  `gorm:"type:text[]" json:"alias"`
	Gender    Gender    `json:"gender"`
	Rlation   string    `json:"relation"`
	Summary   string    `json:"summary"`
	Cover     string    `json:"cover"`
	Images    []string  `gorm:"type:text[]" json:"images"`
	CV        Person    `gorm:"foreignKey:ActorID" json:"cv"`
	ActorID   uint      `gorm:"default:null" json:"-"`
	Tags      []Tag     `gorm:"many2many:character_tag;" json:"tags"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
}

type CharacterTag struct {
	CharacterID uint `gorm:"primaryKey"`
	TagID       uint `gorm:"primaryKey"`
}

func (CharacterTag) TableName() string {
	return "character_tag"
}
