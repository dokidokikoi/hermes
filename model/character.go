package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/datatypes"
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
	BrithYear int     `json:"birth_year,omitempty"`
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
	ID           uint                        `gorm:"primaryKey" json:"id"`
	RelIDs       datatypes.JSONSlice[string] `gorm:"type:jsonb;" json:"rel_ids"`
	Name         string                      `gorm:"type:varchar(255);index:character_name_idx" json:"name"`
	Alias        datatypes.JSONSlice[string] `gorm:"type:jsonb" json:"alias"`
	Gender       Gender                      `gorm:"type:varchar(10)" json:"gender"`
	Summary      string                      `json:"summary"`
	Cover        string                      `gorm:"type:varchar(512)" json:"cover"`
	Images       datatypes.JSONSlice[string] `gorm:"type:jsonb" json:"images"`
	PersonID     uint                        `json:"-"`
	Weight       int8                        `gorm:"default:0;type:int8" json:"weight"`
	PersonalInfo PersonalInfo                `gorm:"type:jsonb" json:"personal_info"`
	CreatedAt    time.Time                   `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt    time.Time                   `gorm:"autoUpdateTime:milli"`

	CV Person `gorm:"foreignKey:PersonID" json:"cv"`
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
