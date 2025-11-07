package model

import (
	"time"
)

type Gender string

const (
	UnKnown Gender = "unknow"
	Male    Gender = "male"
	Female  Gender = "female"
	Futa    Gender = "futa"
)

type Person struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `gorm:"uniqueIndex;type:varchar(100)" json:"uuid"`
	VNDBID    string    `gorm:"type:varchar(255);index:person_vndbid_idx" json:"vndb_id"`
	Name      string    `gorm:"type:varchar(255);index:people_name_idx" json:"name"`
	Alias     Array     `gorm:"type:json" json:"alias"`
	Cover     string    `gorm:"type:varchar(512)" json:"cover"`
	Images    Array     `gorm:"type:json" json:"images"`
	Tags      []Tag     `gorm:"many2many:person_tag;" json:"tags"`
	Summary   string    `json:"summary"`
	Gender    Gender    `gorm:"type:varchar(10)" json:"gender"`
	Links     Links     `gorm:"type:json" json:"links"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
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
