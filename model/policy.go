package model

import (
	"encoding/json"
	"hermes/config"
	"time"
)

type Policy struct {
	ID        uint `gorm:"primaryKey"`
	Policy    string
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
	Proxy config.ProxyConfig `json:"proxy"`
}

type PlatformPolicy []string
type LanguagePolicy []string

type ScraperPolicy struct {
	ScraperName string            `json:"scraper_name"`
	Header      map[string]string `json:"header"`
}
