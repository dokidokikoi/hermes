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
)

var getChuScraper *getchu.GetChu

func init() {
	os.Setenv("https_proxy", "socks://127.0.0.1:7890")
	os.Setenv("http_proxy", "socks://127.0.0.1:7890")
	os.Setenv("all_proxy", "socks://127.0.0.1:7890")
	scraper := getchu.NewGetChu(map[string]string{
		"User-Agent":         config.DefaultUserAgent,
		"Accept-Language":    config.ZhLanguage,
		"Cookie":             getchu.DefaultHeader_Cookie,
		"Sec-Ch-Ua":          getchu.DefaultHeader_SecChUa,
		"Sec-Ch-Ua-Mobile":   getchu.DefaultHeader_SecChUaMobile,
		"Sec-Ch-Ua-Platform": getchu.DefaultHeader_SecChUaPlatform,
		"Referer":            getchu.DefaultHeader_Referer,
	}, "")
	getChuScraper = scraper.(*getchu.GetChu)
}

func TestSearch(t *testing.T) {
	items, err := getChuScraper.Search("彼女", 1)
	if err != nil {
		panic(err)
	}
	for _, item := range items {
		fmt.Printf("%+v\n", item)
	}
}

func TestGetItem(t *testing.T) {
	item, err := getChuScraper.GetItem("https://www.getchu.com/soft.phtml?id=1273918")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", item)
}

func TestGetItemName(t *testing.T) {
	data, err := getChuScraper.DoReq(http.MethodGet, "https://www.getchu.com/soft.phtml?id=1282568", nil, nil)
	if err != nil {
		panic(err)
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	name, err := getChuScraper.GetItemName(root)
	if err != nil {
		panic(err)
	}
	fmt.Println(name)
}

func TestGetItemCover(t *testing.T) {
	data, err := getChuScraper.DoReq(http.MethodGet, "https://www.getchu.com/soft.phtml?id=1273918", nil, nil)
	if err != nil {
		panic(err)
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	cover, images, err := getChuScraper.GetItemCover(root, "1273918")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s, %+v\n", cover, images)
}

func TestGetItemCharacter(t *testing.T) {
	data, err := getChuScraper.DoReq(http.MethodGet, "https://www.getchu.com/soft.phtml?id=1273918", nil, nil)
	if err != nil {
		panic(err)
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	cs, err := getChuScraper.GetItemCharacter(root)
	if err != nil {
		panic(err)
	}
	for _, c := range cs {
		fmt.Printf("%+v\n", c)
	}
}

func TestGetItemStory(t *testing.T) {
	data, err := getChuScraper.DoReq(http.MethodGet, "https://www.getchu.com/soft.phtml?id=1273918", nil, nil)
	if err != nil {
		panic(err)
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	story, err := getChuScraper.GetItemStory(root)
	if err != nil {
		panic(err)
	}

	fmt.Println(story)
}

func TestReq(t *testing.T) {
	data, err := getChuScraper.DoReq(http.MethodGet, "https://www.getchu.com/php/search.phtml?search_keyword=&list_count=30&sort=sales&sort2=down&search_title=%C8%E0%BD%F7&search_brand=&search_person=&search_jan=&search_isbn=&genre=pc_soft&start_date=&end_date=&age=&list_type=list&search=search&pageID=1", nil, nil)
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
		return getChuScraper.DoReq(http.MethodGet, url, nil, nil)
	}, []string{"https://www.getchu.com/brandnew/1273918/c1273918sample3.jpg"})

	fmt.Printf("%+v\n", res)
}
