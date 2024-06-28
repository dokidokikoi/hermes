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
	bangumi.BangumiScraper.GetName(): bangumi.BangumiScraper,
	dlsite.DlSiteScraper.GetName():   dlsite.DlSiteScraper,
	getchu.GetChuScraper.GetName():   getchu.GetChuScraper,
	ggbases.GGBasesScraper.GetName(): ggbases.GGBasesScraper,
	twodfan.TwoDFanScraper.GetName(): twodfan.TwoDFanScraper,
}

var GameScraperPolicyMap = model.ScraperPolicy{
	bangumi.BangumiScraper.GetName(): model.ScraperSubPolicy{
		Header: map[string]string{
			"User-Agent":    "dokidokikoi/meta-scraper (https://github.com/dokidokikoi/meta-scraper)",
			"Authorization": "Bearer ",
		},
	},
	dlsite.DlSiteScraper.GetName(): model.ScraperSubPolicy{
		Header: map[string]string{
			"Sec-Ch-Ua":          `Google Chrome";v="125", "Chromium";v="125", "Not.A/Brand";v="24"`,
			"Sec-Ch-Ua-Mobile":   "?0",
			"Sec-Ch-Ua-Platform": "macOS",
			"User-Agent":         config.DefaultUserAgent,
			"Referer":            "https://www.getchu.com/php/search.phtml?search_keyword=%C8%E0%BD%F7&list_count=30&sort=sales&sort2=down&search_title=&search_brand=&search_person=&search_jan=&search_isbn=&genre=pc_soft&start_date=&end_date=&age=&list_type=list&search=1&pageID=1",
			"Accept-Language":    config.ZhLanguage,
			"Cookie":             "__DLsite_SID=782pmg62psm037ve711d3mcvou; _vwo_uuid_v2=D58EED3043B8C8712836CA3A0CEE347EA|60e7665c582ba675b0d38e5d2fff3d4a; _gcl_au=1.1.1614214747.1716352480; uniqid=0.1jznx3ayl8r; _inflow_ad_params=%7B%22ad_name%22%3A%22organic%22%7D; _fbp=fb.1.1716352481206.123734068; _gaid=876588495.1716352481; _yjsu_yjad=1716352481.85b8a4d8-d0ef-40c3-9add-72cdbb9aefb2; __lt__cid=d543a8ce-a077-4fc3-9c2e-8f3174d95c59; localesuggested=true; locale=zh-cn; _tt_enable_cookie=1; _ttp=rjrtRsh6ouv0PmOqGCRhvWVWaIs; _im_vid=01HYF98ZMNGPY3Z39Q5YJMZW7W; universe_aid=bcf505d16a92b2a620515be740e116240a1a00eccd6e9b0e; adr_id=S7YehhRFnRk3O6gLUtpysCXLJ0EzAGzX1yWf6W6kB4FhT3yt; adultchecked=1; _inflow_params=%7B%22referrer_uri%22%3A%22www.google.com.hk%22%7D; _gid=GA1.2.771398923.1717239817; _ga_QEETZHFB1S=GS1.1.1717290605.1.1.1717290605.0.0.0; _ga_YG879NVEC7=GS1.1.1717290602.1.1.1717290636.0.0.0; _ga_sid=1717304107; __lt__sid=a4160a04-dddad713; DL_PRODUCT_LOG=%2CVJ011538%2CVJ01001190%2CVJ01001393%2CVJ01002056; OptanonConsent=isGpcEnabled=0&datestamp=Sun+Jun+02+2024+14%3A08%3A57+GMT%2B0800+(%E4%B8%AD%E5%9B%BD%E6%A0%87%E5%87%86%E6%97%B6%E9%97%B4)&version=6.23.0&isIABGlobal=false&hosts=&landingPath=NotLandingPage&groups=C0001%3A1%2CC0002%3A0%2CC0003%3A0%2CC0004%3A0&AwaitingReconsent=false; _ga_ZW5GTXK6EV=GS1.1.1717304107.5.1.1717308538.0.0.0; _ga=GA1.1.876588495.1716352481; _inflow_dlsite_params=%7B%22dlsite_referrer_url%22%3A%22https%3A%2F%2Fwww.dlsite.com%2Fpro%2Fwork%2F%3D%2Fproduct_id%2FVJ01001190.html%22%7D",
		},
	},
	getchu.GetChuScraper.GetName(): model.ScraperSubPolicy{
		Header: map[string]string{
			"Sec-Ch-Ua":          `Google Chrome";v="125", "Chromium";v="125", "Not.A/Brand";v="24"`,
			"Sec-Ch-Ua-Mobile":   "?0",
			"Sec-Ch-Ua-Platform": "macOS",
			"User-Agent":         config.DefaultUserAgent,
			"Referer":            "https://www.getchu.com/php/search.phtml?search_keyword=%C8%E0%BD%F7&list_count=30&sort=sales&sort2=down&search_title=&search_brand=&search_person=&search_jan=&search_isbn=&genre=pc_soft&start_date=&end_date=&age=&list_type=list&search=1&pageID=1",
			"Accept-Language":    config.ZhLanguage,
			"Cookie":             "_im_vid=01HYF9KCRA1MT8HSM4EWETGX8S; _gid=GA1.2.1781699859.1717215574; getchu_adalt_flag=getchu.com; ITEM_HISTORY=1282568%7C1273918; _ga_BSNR8334HV=GS1.1.1717222828.5.1.1717225315.53.0.0; _ga_JBMY6G3QFS=GS1.1.1717222828.5.1.1717225315.53.0.0; _ga=GA1.2.1343565952.1716352800; _gat=1",
		},
	},
	ggbases.GGBasesScraper.GetName(): model.ScraperSubPolicy{
		Header: map[string]string{
			"User-Agent":      config.DefaultUserAgent,
			"Referer":         ggbases.GGBasesDomain,
			"Accept-Language": config.ZhLanguage,
		},
	},
	twodfan.TwoDFanScraper.GetName(): model.ScraperSubPolicy{
		Header: map[string]string{
			"User-Agent":      config.DefaultUserAgent,
			"Referer":         "https://2dfan.com/",
			"Accept-Language": config.ZhLanguage,
			"Cookie":          "_ga=GA1.1.566177421.1716285606; pop-blocked=true; _project_hgc_session=amhvTGpZYTdmc3VidU4yQUc2cm01aFdTQzhlTk9NdjI2MXFRVWFsUUw3dmRLTXZ4blYwZ2Q4ZUFOOGtkMld2aTg2YWFtSEpzOFJjTkZSejMvaXg5UytTVzYramdaNzNzbFRXYXJ6a1VLNW5RRzU1L29TK3lyWWJaY0wyVWFKUnN2UDQ0K0hPV2ZDTWx0UFVLdE1tajZ6QndtOGRnWkRndFZIM3BkR0FmaUxVWG5PeGtaeEczRXVWTngvd2hQY25MLS1EbzhJc1ZsbFp3VS92dy8wWGIwWG1nPT0%3D--68accc4aae207d572af489e2c4cfa260efdd5f57; _ga_RF77TZ6QMN=GS1.1.1716638766.7.1.1716641708.0.0.0",
		},
	},
}
