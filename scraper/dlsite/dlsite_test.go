package dlsite_test

import (
	"bytes"
	"fmt"
	"hermes/config"
	"hermes/scraper/dlsite"
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

func TestSearch(t *testing.T) {
	items, err := dlsite.DlSiteScraper.Search("彼女", 1)
	if err != nil {
		panic(err)
	}
	for _, item := range items {
		fmt.Printf("%s ---- %s\n", item.Name, item.Cover)
	}
}

func TestGetItemCover(t *testing.T) {
	data, err := dlsite.DlSiteScraper.DoReq(http.MethodGet, "https://www.dlsite.com/pro/work/=/product_id/VJ01001190.html", nil, nil)
	if err != nil {
		panic(err)
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	cover, images, err := dlsite.DlSiteScraper.GetItemCover(root)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s, %+v", cover, images)
}

func TestGetItemCharacter(t *testing.T) {
	data, err := dlsite.DlSiteScraper.DoReq(http.MethodGet, "https://www.dlsite.com/pro/work/=/product_id/VJ01001190.html", nil, nil)
	if err != nil {
		panic(err)
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	characters, err := dlsite.DlSiteScraper.GetItemCharacter(root)
	if err != nil {
		panic(err)
	}
	for _, c := range characters {
		fmt.Printf("%+v\n", c.CV)
	}
}

func TestGetItemStory(t *testing.T) {
	data, err := dlsite.DlSiteScraper.DoReq(http.MethodGet, "https://www.dlsite.com/pro/work/=/product_id/VJ01001190.html", nil, nil)
	if err != nil {
		panic(err)
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	story, err := dlsite.DlSiteScraper.GetItemStory(root)
	if err != nil {
		panic(err)
	}
	fmt.Println(story)
}

func TestGetItemlinks(t *testing.T) {
	data, err := dlsite.DlSiteScraper.DoReq(http.MethodGet, "https://www.dlsite.com/pro/work/=/product_id/VJ01001190.html", nil, nil)
	if err != nil {
		panic(err)
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	links, err := dlsite.DlSiteScraper.GetItemlink(root, "VJ01001190")
	if err != nil {
		panic(err)
	}
	for _, l := range links {
		fmt.Printf("%+v", l)
	}
}

func TestGetItem(t *testing.T) {
	item, err := dlsite.DlSiteScraper.GetItem("https://www.dlsite.com/pro/work/=/product_id/VJ01001190.html")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", item)
}

func TestReq(t *testing.T) {
	data, err := dlsite.DlSiteScraper.DoReq(http.MethodGet, "https://www.dlsite.com/pro/work/=/product_id/VJ01001190.html", nil, nil)
	if err != nil {
		panic(err)
	}
	f, err := os.Create("index.html")
	if err != nil {
		panic(err)
	}
	fmt.Fprint(f, string(data))
}
