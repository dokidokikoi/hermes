package twodfan_test

import (
	"bytes"
	"fmt"
	"hermes/config"
	"hermes/scraper/twodfan"
	"net/http"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"resty.dev/v3"
)

func init() {
	config.SetConfig("../../conf/application.yaml")
	zaplog.SetLogger(config.GetConfig().LogConfig)
}

func TestTwoDFan_GetItem(t *testing.T) {
	item, err := twodfan.TwoDFanScraper.GetItem("https://2dfan.com/subjects/4566")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", item)
}

func TestTwoDFan_Search(t *testing.T) {
	items, err := twodfan.TwoDFanScraper.Search("ボクの彼女", 1)
	if err != nil {
		panic(err)
	}
	for _, item := range items {
		fmt.Printf("%+v\n", item.Cover)
	}

	fmt.Println(len(items))
}

func TestMain(t *testing.T) {
	os.Setenv("https_proxy", "socks5://127.0.0.1:20170")
	os.Setenv("http_proxy", "socks5://127.0.0.1:20170")
	os.Setenv("all_proxy", "socks5://127.0.0.1:20170")
	rsp, err := resty.New().R().SetHeaders(map[string]string{
		"user-agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36",
		"cookie":     "_ga=GA1.1.274316804.1751095279; _project_hgc_session=T1t5cqvAv1u%2FHORlI7WY3iq69cwcbsgTUQofgmhUv3h4wuvte1P3zH6sDzzj6hL%2FBI8A5k27FvMeIIC%2FUCr7Hdv%2BVRuSHIqFXIWnT6EiNCZZ2xX95PdtRcPdiWNTiQZoXPReJXTP1wdtAt%2FvPiVTzUGbW8ljZdZiLjBuCodtao%2Fz801mdfXUEc3CWBYf5naLzsGHMoRgiBS6Xm38MiQIKZpkPtkZJ3nIGHULivZpt%2Bck15Y7AxvrJgcYK%2FGaNnrVrX6HyPzE2ZutudLWA1WSgUWiEbaGaumHEK3HvwRxkcVpKa%2BhwLZf%2BJujTa7fdFs%2FSabDHJRdQGeiEI8SwvHG3F0%2BgKRCsWIKmshvKDTjX3%2FLGOfNBzDNf2b7N2I%3D--NxMi4KzkxmkoBOoZ--MoOeff9WMBKhr1nDwHjF3g%3D%3D; _ga_RF77TZ6QMN=GS2.1.s1751095279$o1$g1$t1751096558$j60$l0$h0",
	}).Get("https://2dfan.com/subjects/search?keyword=%E3%83%9C%E3%82%AF%E3%81%AE%E5%BD%BC%E5%A5%B3")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(rsp.Bytes()))
}

func TestTwoDFan_GetItemName(t *testing.T) {
	items, err := twodfan.TwoDFanScraper.Search("ボクの彼女", 1)
	if err != nil {
		panic(err)
	}

	if len(items) > 0 {
		url := twodfan.TwoDFanScraper.AbsUrl(items[0].URl)
		data, err := twodfan.TwoDFanScraper.DoReq(http.MethodGet, url, nil, nil)
		if err != nil {
			panic(err)
		}

		root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
		if err != nil {
			panic(err)
		}

		name, alias, err := twodfan.TwoDFanScraper.GetItemName(root)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s, %+v\n", name, alias)
	}
}

func TestTwoDFan_GetItemCover(t *testing.T) {
	items, err := twodfan.TwoDFanScraper.Search("ボクの彼女", 1)
	if err != nil {
		panic(err)
	}

	if len(items) > 0 {
		url := twodfan.TwoDFanScraper.AbsUrl(items[0].URl)
		data, err := twodfan.TwoDFanScraper.DoReq(http.MethodGet, url, nil, nil)
		if err != nil {
			panic(err)
		}

		root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
		if err != nil {
			panic(err)
		}

		cover, images, err := twodfan.TwoDFanScraper.GetItemCover(root)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s, %+v\n", cover, images)
	}
}

func TestTwoDFan_GetItemLinks(t *testing.T) {
	items, err := twodfan.TwoDFanScraper.Search("ボクの彼女はガテン系", 1)
	if err != nil {
		panic(err)
	}

	if len(items) > 0 {
		url := twodfan.TwoDFanScraper.AbsUrl(items[0].URl)
		data, err := twodfan.TwoDFanScraper.DoReq(http.MethodGet, url, nil, nil)
		if err != nil {
			panic(err)
		}

		root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
		if err != nil {
			panic(err)
		}

		links, err := twodfan.TwoDFanScraper.GetItemLinks(root)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", links)
	}
}

func TestTwoDFan_GetItemIssueDate(t *testing.T) {
	items, err := twodfan.TwoDFanScraper.Search("ボクの彼女はガテン系", 1)
	if err != nil {
		panic(err)
	}

	if len(items) > 0 {
		url := twodfan.TwoDFanScraper.AbsUrl(items[0].URl)
		data, err := twodfan.TwoDFanScraper.DoReq(http.MethodGet, url, nil, nil)
		if err != nil {
			panic(err)
		}

		root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
		if err != nil {
			panic(err)
		}

		t, err := twodfan.TwoDFanScraper.GetItemIssueDate(root)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", t.Format("2006-01-02"))
	}
}

func TestTwoDFan_GetItemTags(t *testing.T) {
	items, err := twodfan.TwoDFanScraper.Search("ボクの彼女はガテン系", 1)
	if err != nil {
		panic(err)
	}

	if len(items) > 0 {
		url := twodfan.TwoDFanScraper.AbsUrl(items[0].URl)
		data, err := twodfan.TwoDFanScraper.DoReq(http.MethodGet, url, nil, nil)
		if err != nil {
			panic(err)
		}

		root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
		if err != nil {
			panic(err)
		}

		tags, err := twodfan.TwoDFanScraper.GetItemTags(root)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", tags)
	}
}

func TestTwoDFan_GetItemStaff(t *testing.T) {
	items, err := twodfan.TwoDFanScraper.Search("ボクの彼女はガテン系", 1)
	if err != nil {
		panic(err)
	}

	if len(items) > 0 {
		url := twodfan.TwoDFanScraper.AbsUrl(items[0].URl)
		data, err := twodfan.TwoDFanScraper.DoReq(http.MethodGet, url, nil, nil)
		if err != nil {
			panic(err)
		}

		root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
		if err != nil {
			panic(err)
		}

		staff, err := twodfan.TwoDFanScraper.GetItemStaff(root)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", staff)
	}
}

func TestTwoDFan_GetItemStory(t *testing.T) {
	items, err := twodfan.TwoDFanScraper.Search("ボクの彼女はガテン系", 1)
	if err != nil {
		panic(err)
	}

	if len(items) > 0 {
		url := twodfan.TwoDFanScraper.AbsUrl(items[0].URl)
		data, err := twodfan.TwoDFanScraper.DoReq(http.MethodGet, url, nil, nil)
		if err != nil {
			panic(err)
		}

		root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
		if err != nil {
			panic(err)
		}

		story, images, err := twodfan.TwoDFanScraper.GetItemStory(root)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s, %+v\n", story, images)
	}
}

func TestTwoDFan_GetItemPublisher(t *testing.T) {
	url := "https://2dfan.com/subjects/4566"
	data, err := twodfan.TwoDFanScraper.DoReq(http.MethodGet, url, nil, nil)
	if err != nil {
		panic(err)
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	publisher, err := twodfan.TwoDFanScraper.GetItemPublisher(root)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", publisher)
}

func TestReq(t *testing.T) {
	data, err := twodfan.TwoDFanScraper.DoReq(http.MethodGet, "https://2dfan.com/subjects/2020", nil, nil)
	if err != nil {
		panic(err)
	}
	f, err := os.Create("index.html")
	if err != nil {
		panic(err)
	}
	fmt.Fprint(f, string(data))
}
