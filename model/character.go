package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Trait struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type PersonalInfo struct {
	BloodType string  `json:"blood_type,omitempty"`
	Height    int     `json:"height,omitempty"`
	Weight    int     `json:"weight,omitempty"`
	Bust      int     `json:"bust,omitempty"`
	Waist     int     `json:"waist,omitempty"`
	Hips      int     `json:"hips,omitempty"`
	Cup       string  `json:"cup,omitempty"`
	Age       int     `json:"age,omitempty"`
	Birthday  [2]int  `json:"birthday,omitempty"`
	Traits    []Trait `json:"traits,omitempty"`
}

func (p *PersonalInfo) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, p)
	case string:
		return json.Unmarshal([]byte(src), p)
	case nil:
		*p = PersonalInfo{}
		return nil
	}

	return fmt.Errorf("cannot convert %T to Link", src)
}

func (p PersonalInfo) Value() (driver.Value, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return "{}", err
	}

	return data, nil
}

type Character struct {
	ID           uint         `gorm:"primaryKey" json:"id"`
	UUID         string       `gorm:"uniqueIndex;type:varchar(255)" json:"uuid"`
	VNDBID       string       `gorm:"type:varchar(255);index:char_vndbid_idx" json:"vndb_id"`
	Name         string       `gorm:"type:varchar(255);index:character_name_idx" json:"name"`
	Alias        Array        `gorm:"type:json" json:"alias"`
	Gender       Gender       `gorm:"type:varchar(10)" json:"gender"`
	Summary      string       `json:"summary"`
	Cover        string       `gorm:"type:varchar(512)" json:"cover"`
	Images       Array        `gorm:"type:json" json:"images"`
	CV           Person       `gorm:"foreignKey:PersonID" json:"cv"`
	PersonID     uint         `gorm:"default:0" json:"-"`
	Tags         []Tag        `gorm:"many2many:character_tag;" json:"tags"`
	Weight       int8         `gorm:"default:0;type:tinyint" json:"weight"`
	PersonalInfo PersonalInfo `json:"personal_info"`
	CreatedAt    time.Time    `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt    time.Time    `gorm:"autoUpdateTime:milli"`
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
