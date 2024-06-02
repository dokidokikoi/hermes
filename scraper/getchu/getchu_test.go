package getchu

import (
	"bytes"
	"fmt"
	"hermes/config"
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
	items, err := GetChuScraper.Sreach("彼女", 1)
	if err != nil {
		panic(err)
	}
	for _, item := range items {
		fmt.Printf("%+v\n", item)
	}
}

func TestGetItem(t *testing.T) {
	item, err := GetChuScraper.GetItem("https://www.getchu.com/soft.phtml?id=1273918")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", item)
}

func TestGetItemName(t *testing.T) {
	data, err := GetChuScraper.DoReq(http.MethodGet, "https://www.getchu.com/soft.phtml?id=1282568", nil, nil)
	if err != nil {
		panic(err)
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	name, err := GetChuScraper.GetItemName(root)
	if err != nil {
		panic(err)
	}
	fmt.Println(name)
}

func TestGetItemCover(t *testing.T) {
	data, err := GetChuScraper.DoReq(http.MethodGet, "https://www.getchu.com/soft.phtml?id=1273918", nil, nil)
	if err != nil {
		panic(err)
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	cover, images, err := GetChuScraper.GetItemCover(root, "1273918")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s, %+v\n", cover, images)
}

func TestGetItemCharacter(t *testing.T) {
	data, err := GetChuScraper.DoReq(http.MethodGet, "https://www.getchu.com/soft.phtml?id=1273918", nil, nil)
	if err != nil {
		panic(err)
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	cs, err := GetChuScraper.GetItemCharacter(root)
	if err != nil {
		panic(err)
	}
	for _, c := range cs {
		fmt.Printf("%+v\n", c)
	}
}

func TestGetItemStory(t *testing.T) {
	data, err := GetChuScraper.DoReq(http.MethodGet, "https://www.getchu.com/soft.phtml?id=1273918", nil, nil)
	if err != nil {
		panic(err)
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	story, err := GetChuScraper.GetItemStory(root)
	if err != nil {
		panic(err)
	}

	fmt.Println(story)
}

func TestReq(t *testing.T) {
	data, err := GetChuScraper.DoReq(http.MethodGet, "https://www.getchu.com/soft.phtml?id=1282568", nil, nil)
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
		return GetChuScraper.DoReq(http.MethodGet, url, nil, nil)
	}, []string{"https://www.getchu.com/brandnew/1273918/c1273918sample3.jpg"})

	fmt.Printf("%+v\n", res)
}
