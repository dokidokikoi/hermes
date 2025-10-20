package model

import (
	"encoding/json"
	"time"
)

type Policy struct {
	ID        uint      `gorm:"primaryKey;type:bigint"`
	Key       string    `gorm:"type:varchar(32)"`
	Policy    string    `gorm:"type:json"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
}

func (Policy) TableName() string {
	return "policies"
}

func Parse[T SystemPolicy | PlatformPolicy | LanguagePolicy | ScraperPolicy](str string) (*T, error) {
	t := new(T)
	err := json.Unmarshal([]byte(str), t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

type SystemPolicy struct {
	GameLibrary string `json:"game_library"`
}

func (SystemPolicy) Key() string {
	return "system"
}

type PlatformPolicy []string

func (PlatformPolicy) Key() string {
	return "platform"
}

type LanguagePolicy []string

func (LanguagePolicy) Key() string {
	return "language"
}

type ScraperPolicy map[string]ScraperSubPolicy
type ScraperSubPolicy struct {
	Header map[string]string `json:"header"`
}

func (ScraperPolicy) Key() string {
	return "scraper"
}
