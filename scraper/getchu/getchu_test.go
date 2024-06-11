package getchu_test

import (
	"bytes"
	"fmt"
	"hermes/config"
	"hermes/scraper/getchu"
	"hermes/tools"
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
	items, err := getchu.GetChuScraper.Search("彼女", 1)
	if err != nil {
		panic(err)
	}
	for _, item := range items {
		fmt.Printf("%+v\n", item)
	}
}

func TestGetItem(t *testing.T) {
	item, err := getchu.GetChuScraper.GetItem("https://www.getchu.com/soft.phtml?id=1273918")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", item)
}

func TestGetItemName(t *testing.T) {
	data, err := getchu.GetChuScraper.DoReq(http.MethodGet, "https://www.getchu.com/soft.phtml?id=1282568", nil, nil)
	if err != nil {
		panic(err)
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	name, err := getchu.GetChuScraper.GetItemName(root)
	if err != nil {
		panic(err)
	}
	fmt.Println(name)
}

func TestGetItemCover(t *testing.T) {
	data, err := getchu.GetChuScraper.DoReq(http.MethodGet, "https://www.getchu.com/soft.phtml?id=1273918", nil, nil)
	if err != nil {
		panic(err)
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	cover, images, err := getchu.GetChuScraper.GetItemCover(root, "1273918")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s, %+v\n", cover, images)
}

func TestGetItemCharacter(t *testing.T) {
	data, err := getchu.GetChuScraper.DoReq(http.MethodGet, "https://www.getchu.com/soft.phtml?id=1273918", nil, nil)
	if err != nil {
		panic(err)
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	cs, err := getchu.GetChuScraper.GetItemCharacter(root)
	if err != nil {
		panic(err)
	}
	for _, c := range cs {
		fmt.Printf("%+v\n", c)
	}
}

func TestGetItemStory(t *testing.T) {
	data, err := getchu.GetChuScraper.DoReq(http.MethodGet, "https://www.getchu.com/soft.phtml?id=1273918", nil, nil)
	if err != nil {
		panic(err)
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	story, err := getchu.GetChuScraper.GetItemStory(root)
	if err != nil {
		panic(err)
	}

	fmt.Println(story)
}

func TestReq(t *testing.T) {
	data, err := getchu.GetChuScraper.DoReq(http.MethodGet, "https://www.getchu.com/php/search.phtml?search_keyword=&list_count=30&sort=sales&sort2=down&search_title=%C8%E0%BD%F7&search_brand=&search_person=&search_jan=&search_isbn=&genre=pc_soft&start_date=&end_date=&age=&list_type=list&search=search&pageID=1", nil, nil)
	if err != nil {
		panic(err)
	}
	f, err := os.Create("index.html")
	if err != nil {
		panic(err)
	}
	html, _ := tools.Jp2Utf8(data)
	fmt.Fprint(f, string(html))
}

func TestGetImage(t *testing.T) {
	res := tools.SaveBunchTmpFile(func(url string) ([]byte, error) {
		return getchu.GetChuScraper.DoReq(http.MethodGet, url, nil, nil)
	}, []string{"https://www.getchu.com/brandnew/1273918/c1273918sample3.jpg"})

	fmt.Printf("%+v\n", res)
}
