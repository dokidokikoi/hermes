package ggbases

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hermes/config"
	"hermes/internal/handler"
	"hermes/model"
	"hermes/scraper"
	"hermes/tools"
	"io"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	comm_tools "github.com/dokidokikoi/go-common/tools"
)

var GGBasesDomain = "https://ggbases.dlgal.com/"
var GGBasesSearchUri = "https://ggbases.dlgal.com/search.so?p=0&title=%d&advanced=0"
var GGBasesMagnetUri = "https://ggbases.dlgal.com/magnet.so?id=%s"
var GGBasesBtUri = "https://ggbases.dlgal.com/down.so?id=%s"

type GGBases struct {
	Name      string
	Domain    string
	SearchUri string
	Headers   map[string]string
}

var GGBasesScraper *GGBases

func (gg *GGBases) DoReq(method, uri string, header map[string]string, body interface{}) ([]byte, error) {
	h := map[string]string{}
	for k, v := range gg.Headers {
		h[k] = v
	}
	for k, v := range header {
		h[k] = v
	}

	var r io.Reader
	if method == http.MethodGet {
		query := comm_tools.GenQueryParams(body)
		uri += "?" + query
	} else {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		r = bytes.NewBuffer(data)
	}

	data, _, err := tools.MakeRequest(method, uri, config.GetConfig().ProxyConfig, r, h, nil)
	return data, err
}

// func (gg *GGBases) DoChromeReq(url string, headless bool, fs ...func(ctx context.Context)) ([]byte, error) {
// 	opts := append(chromedp.DefaultExecAllocatorOptions[:],
// 		chromedp.Flag("headless", headless),
// 		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36"),
// 	)
// 	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
// 	defer cancel()
// 	chromeCtx, cancel := chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
// 	defer cancel()
// 	// 保持浏览器窗口开启
// 	_ = chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

// 	timeOutCtx, cancel := context.WithTimeout(chromeCtx, 60*time.Second)
// 	defer cancel()

// 	var htmlContent string
// 	err := chromedp.Run(timeOutCtx,
// 		network.Enable(),
// 		//需要爬取的网页的url
// 		chromedp.Navigate(url),
// 		network.SetExtraHTTPHeaders(map[string]interface{}{"Accept-Language": "zh-cn,zh;q=0.5", "X-Forwarded-For": "https://ggbases.dlgal.com/"}),
// 		chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery),
// 	)
// 	for _, f := range fs {
// 		f(timeOutCtx)
// 	}
// 	return []byte(htmlContent), err
// }

func (gg *GGBases) GetItem(uri string) (*scraper.GameItem, error) {
	data, err := gg.DoReq(http.MethodGet, uri, nil, nil)
	if err != nil {
		return nil, err
	}
	item := &scraper.GameItem{GameVo: handler.GameVo{Links: []model.Link{{Name: "ggbases", Url: uri}}}}
	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (gg GGBases) GetItemName(node *goquery.Document) (string, error) {
	return node.Find("#atitle").Text(), nil
}

func (gg GGBases) GetItemCover(node *goquery.Document) ([]string, error) {

	return nil, nil
}

func (gg GGBases) GetItemTags(node *goquery.Document) ([]model.Tag, error) {

	return nil, nil
}

func (gg GGBases) GetItemIssueeDate(node *goquery.Document) (time.Time, error) {
	return time.Time{}, nil
}
func (gg GGBases) GetItemLink(node *goquery.Document) (string, error) {
	return "", nil
}
func (gg GGBases) GetItemInformation(node *goquery.Document) ([]string, error) {
	return nil, nil
}
func (gg GGBases) GetItemSaveData(node *goquery.Document) (string, error) {
	link, _ := node.Find("#touch tbody>tr:nth-child(7)>td a:nth-child(2)").Attr("href")
	return fmt.Sprintf("https:%s", link), nil
}
func (gg GGBases) GetItemWalkThrough(node *goquery.Document) (string, error) {
	link, _ := node.Find("#touch tbody>tr:nth-child(7)>td a:nth-child(3)").Attr("href")
	return fmt.Sprintf("https:%s", link), nil
}
func (gg GGBases) GetItemSize(node *goquery.Document) (string, error) {
	return node.Find("#touch tbody tr:nth-child(5) td:nth-child(2) span").Text(), nil
}

func (gg GGBases) GetItemOtherInfo(node *goquery.Document) (string, error) {
	return node.Find("#description div[markdown-text]").Html()
}

func init() {
	headers := make(map[string]string)
	headers["User-Agent"] = config.DefaultUserAgent
	headers["Referer"] = GGBasesDomain
	headers["Accept-Language"] = config.ZhLanguage
	GGBasesScraper = &GGBases{
		Name:      "ggbases",
		Domain:    GGBasesDomain,
		SearchUri: GGBasesSearchUri,
		Headers:   headers,
	}
}
