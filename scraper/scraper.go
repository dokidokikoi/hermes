package scraper

import (
	"izumi/internal/handler"
)

type GameItem struct {
	handler.GameVo
	ScraperName string `json:"scraper_name"`
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
	GetGameItem(url string) (*GameItem, error)
	SearchGame(keyword string, page int) ([]*SearchItem, error)

	GetName() string
	SetHeader(header map[string]string)
	SetProxy(proxy string)
}

type ComicItem struct {
	handler.GameVo
	ScraperName string `json:"scraper_name"`
}

type IComicScraper interface {
	GetComicItem(url string) (*GameItem, error)
	SearchComic(keyword string, page int) ([]*SearchItem, error)

	GetName() string
	SetHeader(header map[string]string)
	SetProxy(proxy string)
}
