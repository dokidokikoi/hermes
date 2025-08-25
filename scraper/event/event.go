package event

import (
	"hermes/config"
	"hermes/model"
	"hermes/scraper"
	"hermes/scraper/bangumi"
	"hermes/scraper/dlsite"
	"hermes/scraper/getchu"
	"hermes/scraper/ggbases"
	"hermes/scraper/twodfan"
)

var GameScraperMap = map[string]scraper.IGameScraper{
	// bangumi.BangumiScraper.GetName(): bangumi.BangumiScraper,
	// dlsite.DlSiteScraper.GetName():   dlsite.DlSiteScraper,
	// getchu.GetChuScraper.GetName():   getchu.GetChuScraper,
	// ggbases.GGBasesScraper.GetName(): ggbases.GGBasesScraper,
	// twodfan.TwoDFanScraper.GetName(): twodfan.TwoDFanScraper,
}

func RegisterScraper(scraper scraper.IGameScraper) {
	GameScraperMap[scraper.GetName()] = scraper
}

var GameScraperPolicyMap = model.ScraperPolicy{
	bangumi.Name: model.ScraperSubPolicy{
		Header: map[string]string{
			"User-Agent":    bangumi.DefaultHeader_UserAgent,
			"Authorization": "Bearer ",
		},
	},
	dlsite.Name: model.ScraperSubPolicy{
		Header: map[string]string{
			"Sec-Ch-Ua":          dlsite.DefaultHeader_SecChUa,
			"Sec-Ch-Ua-Mobile":   dlsite.DefaultHeader_SecChUaMobile,
			"Sec-Ch-Ua-Platform": dlsite.DefaultHeader_SecChUaPlatform,
			"User-Agent":         config.DefaultUserAgent,
			"Accept-Language":    config.ZhLanguage,
			"Cookie":             dlsite.DefaultHeader_Cookie,
		},
	},
	getchu.Name: model.ScraperSubPolicy{
		Header: map[string]string{
			"Sec-Ch-Ua":          getchu.DefaultHeader_SecChUa,
			"Sec-Ch-Ua-Mobile":   getchu.DefaultHeader_SecChUaMobile,
			"Sec-Ch-Ua-Platform": getchu.DefaultHeader_SecChUaPlatform,
			"User-Agent":         config.DefaultUserAgent,
			"Referer":            getchu.DefaultHeader_Referer,
			"Accept-Language":    config.ZhLanguage,
			"Cookie":             getchu.DefaultHeader_Cookie,
		},
	},
	ggbases.Name: model.ScraperSubPolicy{
		Header: map[string]string{
			"User-Agent":      config.DefaultUserAgent,
			"Referer":         ggbases.GGBasesDomain,
			"Accept-Language": config.ZhLanguage,
		},
	},
	twodfan.Name: model.ScraperSubPolicy{
		Header: map[string]string{
			"User-Agent":      config.DefaultUserAgent,
			"Referer":         "https://2dfan.com/",
			"Accept-Language": config.ZhLanguage,
			"Cookie":          twodfan.DefaultHeader_Cookie,
		},
	},
}
