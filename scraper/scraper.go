package scraper

import (
	"hermes/internal/handler"
)

type GameItem struct {
	handler.GameVo
	ScraperName string   `json:"scraper_name"`
	AllImages   []string `json:"all_images"`
}

type SearchItem struct {
	Name    string `json:"name"`
	Key     string `json:"key"`
	URl     string `json:"url"`
	Summary string `json:"summary"`
	Cover   string `json:"cover"`

	ScraperName string `json:"scraper_name"`
}

type IGameScraper interface {
	GetItem(url string) (*GameItem, error)
	Search(keyword string, page int) ([]*SearchItem, error)
	GetName() string
	SetHeader(header map[string]string)
}
