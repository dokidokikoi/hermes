package model

import (
	"time"

	"github.com/lib/pq"
)

type Gender string

const (
	UnKnown Gender = "unknow"
	Male    Gender = "male"
	Female  Gender = "female"
	Futa    Gender = "futa"
)

type Person struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name"`
	Alias     pq.StringArray `gorm:"type:text[]" json:"alias"`
	Cover     string         `json:"cover"`
	Images    pq.StringArray `gorm:"type:text[]" json:"images"`
	Tags      []Tag          `gorm:"many2many:person_tag;" json:"tags"`
	Summary   string         `json:"summary"`
	Gender    Gender         `json:"gender"`
	CreatedAt time.Time      `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:milli"`
}

func (Person) TableName() string {
	return "people"
}

type PersonTag struct {
	PersonID uint `gorm:"primaryKey"`
	TagID    uint `gorm:"primaryKey"`
}

func (PersonTag) TableName() string {
	return "person_tag"
}
