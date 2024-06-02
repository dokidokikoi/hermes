package ggbases

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hermes/config"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
)

func init() {
	config.SetConfig("../../conf/application.yaml")
	zaplog.SetLogger(config.GetConfig().LogConfig)
}

func TestSearch(t *testing.T) {
	items, err := GGBasesScraper.Sreach("彼女", 1)
	if err != nil {
		panic(err)
	}
	for _, item := range items {
		fmt.Printf("%+v\n", item)
	}
}

func TestReq(t *testing.T) {
	data, err := GGBasesScraper.DoReq(http.MethodGet, "https://ggbases.dlgal.com/view.so?id=124340", nil, nil)
	if err != nil {
		panic(err)
	}
	str := string(data)
	idx := strings.Index(str, `$("#showCoverBtn").attr("href", "`)
	idx += len(`$("#showCoverBtn").attr("href", "`)

	buf := bytes.Buffer{}
	for ; idx < len(str); idx++ {
		if str[idx] == '"' {
			break
		}
		buf.WriteByte(str[idx])
	}
	fmt.Println(buf.String())

	f, err := os.Create("index.html")
	if err != nil {
		panic(err)
	}

	fmt.Fprint(f, string(data))
}

func TestGetItem(t *testing.T) {
	item, err := GGBasesScraper.GetItem("https://ggbases.dlgal.com/view.so?id=120270")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", item.Staff)
}

func TestGetItemLink(t *testing.T) {
	data, err := GGBasesScraper.DoReq(http.MethodGet, "https://ggbases.dlgal.com/view.so?id=120270", nil, nil)
	if err != nil {
		panic(err)
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	links, err := GGBasesScraper.GetItemLink(root, "120270")
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

	body, err := GGBasesScraper.DoReq(method, url, map[string]string{
		"Content-Type": writer.FormDataContentType(),
		"Referer":      "https://ggbases.dlgal.com/view.so?id=120270",
	}, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
	resp := new(DetailResp)
	err = json.Unmarshal(body, resp)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", resp)
}
