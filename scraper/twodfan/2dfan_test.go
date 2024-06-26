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
