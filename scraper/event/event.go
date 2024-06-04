package event

import (
	"hermes/scraper"
	"hermes/scraper/bangumi"
	"hermes/scraper/dlsite"
	"hermes/scraper/getchu"
	"hermes/scraper/ggbases"
	"hermes/scraper/twodfan"
)

var GameScraperMap = map[string]scraper.IGameScraper{
	bangumi.BangumiScraper.GetName(): bangumi.BangumiScraper,
	dlsite.DlSiteScraper.GetName():   dlsite.DlSiteScraper,
	getchu.GetChuScraper.GetName():   getchu.GetChuScraper,
	ggbases.GGBasesScraper.GetName(): ggbases.GGBasesScraper,
	twodfan.TwoDFanScraper.GetName(): twodfan.TwoDFanScraper,
}
