package ggbases_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hermes/scraper/ggbases"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var ggBasesScraper *ggbases.GGBases

func init() {
	scraper := ggbases.NewGGBases(map[string]string{
		"User-Agent":      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36",
		"Accept-Language": "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		// "Cookie":          ggbases.DefaultHeader_Cookie,
	}, "")
	ggBasesScraper = scraper.(*ggbases.GGBases)
}

func TestSearch(t *testing.T) {
	items, err := ggBasesScraper.Search("彼女", 1)
	if err != nil {
		panic(err)
	}
	for _, item := range items {
		fmt.Printf("%+v\n", item)
	}
}

func TestReq(t *testing.T) {

	data, err := ggBasesScraper.DoReq(http.MethodGet, "https://ggbases.dlgal.com/", nil, nil)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("index.html")
	if err != nil {
		panic(err)
	}

	fmt.Fprint(f, string(data))
}

func TestGetItem(t *testing.T) {
	item, err := ggBasesScraper.GetItem("https://ggbases.dlgal.com/view.so?id=125945")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", item.Staff)
}

func TestGetItemLink(t *testing.T) {
	data, err := ggBasesScraper.DoReq(http.MethodGet, "https://ggbases.dlgal.com/view.so?id=120270", nil, nil)
	if err != nil {
		panic(err)
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	links, err := ggBasesScraper.GetItemLink(root, "120270")
	if err != nil {
		panic(err)
	}

	for _, l := range links {
		fmt.Printf("%+v\n", l)
	}
}

func TestDetail(t *testing.T) {
	url := "https://ggbases.dlgal.com/viewDetail.so?id=1282568"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("id", "120270")
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := ggBasesScraper.DoReq(method, url, map[string]string{
		"Content-Type": writer.FormDataContentType(),
		"Referer":      "https://ggbases.dlgal.com/view.so?id=120270",
	}, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
	resp := new(ggbases.DetailResp)
	err = json.Unmarshal(body, resp)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", resp)
}

func TestGetMagnet(t *testing.T) {
	hash, err := ggBasesScraper.GetMagnet("128413")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hash)
}
