package scraper

import (
	"hermes/internal/handler"
)

type GameItem struct {
	handler.GameVo
	ScraperName string
	AllImages   []string
}

type SearchItem struct {
	Name    string
	Key     string
	URl     string
	Summary string
	Cover   string

	ScraperName string
}

type IGameScraper interface {
	GetItem(url string) (*GameItem, error)
	Search(keyword string, page int) ([]*SearchItem, error)
	GetName() string
	SetHeader(header map[string]string)
}
