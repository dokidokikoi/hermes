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

var GameScraperConstructors = map[string]func(header map[string]string) scraper.IGameScraper{
	bangumi.Name: bangumi.NewBangumi,
	dlsite.Name:  dlsite.NewDlSite,
	getchu.Name:  getchu.NewGetChu,
	ggbases.Name: ggbases.NewGGBases,
	twodfan.Name: twodfan.NewTwoDFan,
}

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
			"Referer":         ggbases.GGBasesDomain,
			"User-Agent":      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36",
			"Accept-Language": "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7",
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
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
