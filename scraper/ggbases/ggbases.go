package ggbases

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hermes/config"
	"hermes/internal/handler"
	"hermes/model"
	"hermes/scraper"
	"hermes/tools"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dokidokikoi/go-common/gopool"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	comm_tools "github.com/dokidokikoi/go-common/tools"
	"github.com/dop251/goja"
	"go.uber.org/zap"
)

var GGBasesDomain = "https://ggbases.dlgal.com/"
var GGBasesSearchUri = "https://ggbases.dlgal.com/search.so?p=%d&title=%s&advanced=0"
var GGBasesMagnetUri = "https://ggbases.dlgal.com/magnet.so?id=%s"
var GGBasesBtUri = "https://ggbases.dlgal.com/down.so?id=%s"
var GGBasesDetailUri = "https://ggbases.dlgal.com/viewDetail.so"

type GGBases struct {
	sync.RWMutex
	name      string
	Domain    string
	SearchUri string
	Headers   map[string]string
}

var GGBasesScraper *GGBases

func init() {
	headers := make(map[string]string)
	headers["User-Agent"] = config.DefaultUserAgent
	headers["Referer"] = GGBasesDomain
	headers["Accept-Language"] = config.ZhLanguage
	GGBasesScraper = &GGBases{
		name:      "ggbases",
		Domain:    GGBasesDomain,
		SearchUri: GGBasesSearchUri,
		Headers:   headers,
	}
}

func (gg *GGBases) GetName() string {
	return gg.name
}

func (gg *GGBases) SetHeader(header map[string]string) {
	gg.Lock()
	for k, v := range header {
		gg.Headers[k] = v
	}
	gg.Unlock()
}

func (gg *GGBases) Search(keyword string, page int) ([]*scraper.SearchItem, error) {
	data, err := gg.DoReq(http.MethodGet, fmt.Sprintf(gg.SearchUri, page-1, keyword), nil, nil)
	if err != nil {
		return nil, err
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	wait := sync.WaitGroup{}
	var lock sync.Mutex
	var items []*scraper.SearchItem
	root.Find("tr.dtr").Each(func(i int, s *goquery.Selection) {
		wait.Add(1)

		gopool.Go(func() {
			defer wait.Done()

			item := &scraper.SearchItem{
				URl:         "https:" + s.Find("td.t-l").Eq(1).Find("a").AttrOr("href", ""),
				Name:        s.Find("td.t-l").Eq(1).Text(),
				ScraperName: gg.name,
			}
			itemData, err := gg.DoReq(http.MethodGet, item.URl, nil, nil)
			if err != nil {
				return
			}
			root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(itemData))
			if err != nil {
				return
			}
			item.Name = root.Find("#atitle").Text()
			item.Summary, _ = root.Find("#description div[markdown-text]").Html()

			item.Cover, err = gg.GetItemCover(root)
			if err != nil {
				zaplog.L().Error("获取封面失败", zap.String("scraper", gg.name), zap.Error(err))
			}

			lock.Lock()
			items = append(items, item)
			lock.Unlock()
		})
	})
	wait.Wait()

	return items, nil
}

func (gg *GGBases) DoReq(method, uri string, header map[string]string, body interface{}) ([]byte, error) {
	h := map[string]string{}
	gg.RLock()
	for k, v := range gg.Headers {
		h[k] = v
	}
	gg.RUnlock()
	for k, v := range header {
		h[k] = v
	}

	var r io.Reader
	if method == http.MethodGet {
		query := comm_tools.GenQueryParams(body)
		if query != "" {
			uri += "?" + query
		}
	} else {
		reader, ok := body.(io.Reader)
		if ok {
			r = reader
		} else {
			data, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			r = bytes.NewBuffer(data)
		}
	}

	data, _, err := tools.MakeRequest(method, uri, config.GetConfig().ProxyConfig, r, h, nil, config.DefaultRetryCnt)
	return data, err
}

func (gg *GGBases) GetItem(uri string) (*scraper.GameItem, error) {
	data, err := gg.DoReq(http.MethodGet, uri, nil, nil)
	if err != nil {
		return nil, err
	}
	item := &scraper.GameItem{GameVo: handler.GameVo{Links: []model.Link{{Name: "ggbases", Url: uri}}}, ScraperName: gg.name}
	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	arr := strings.Split(uri, "?")
	if len(arr) < 2 {
		return nil, errors.New("uri error")
	}
	id := ""
	arr = strings.Split(arr[1], "&")
	for _, a := range arr {
		if a[:3] == "id=" {
			id = a[3:]
			break
		}
	}
	if id == "" {
		return nil, errors.New("uri error")
	}

	item.Name, err = gg.GetItemName(root)
	if err != nil {
		zaplog.L().Error("获取名称失败", zap.String("scraper", gg.name), zap.String("uri", uri), zap.Error(err))
	}
	item.Cover, err = gg.GetItemCover(root)
	if err != nil {
		zaplog.L().Error("获取封面失败", zap.String("scraper", gg.name), zap.String("uri", uri), zap.Error(err))
	}
	item.IssueDate, err = gg.GetItemIssueeDate(root)
	if err != nil {
		zaplog.L().Error("获取发布时间失败", zap.String("scraper", gg.name), zap.String("uri", uri), zap.Error(err))
	}
	links, err := gg.GetItemLink(root, id)
	if err != nil {
		zaplog.L().Error("获取相关链接失败", zap.String("scraper", gg.name), zap.String("uri", uri), zap.Error(err))
	}
	item.Links = append(item.Links, links...)
	item.OtherInfo, err = gg.GetItemOtherInfo(root)
	if err != nil {
		zaplog.L().Error("获取其它信息失败", zap.String("scraper", gg.name), zap.String("uri", uri), zap.Error(err))
	}

	str := string(data)
	idx := strings.Index(str, `var extags = "`)
	idx += len(`var extags = "`)
	if idx != -1 {
		buf := bytes.Buffer{}
		for ; idx < len(str); idx++ {
			if str[idx] == '"' {
				break
			}
			if str[idx] == ';' {
				arr := strings.Split(buf.String(), ":")
				if len(arr) == 2 {
					contents := strings.Split(arr[1], "~")
					content := contents[len(contents)-1]
					switch arr[0] {
					case "artist":
						item.Staff = append(item.Staff, handler.StaffVo{
							Name:     content,
							Relation: []string{model.PRelationPainter.String()},
						})
					case "group":
						item.Publisher = &model.Publisher{
							Name: content,
						}
					default:
						item.Tags = append(item.Tags, model.Tag{
							Name: content,
						})
					}
				}
				buf.Reset()
			}
			buf.WriteByte(str[idx])
		}
	}

	return item, nil
}

func (gg *GGBases) GetItemName(node *goquery.Document) (string, error) {
	return node.Find("#atitle").Text(), nil
}

func (gg *GGBases) GetItemCover(node *goquery.Document) (string, error) {
	html, err := node.Html()
	if err != nil {
		return "", err
	}
	idx := strings.Index(html, `$("#showCoverBtn").attr("href", "`)
	idx += len(`$("#showCoverBtn").attr("href", "`)
	buf := bytes.Buffer{}
	for ; idx < len(html); idx++ {
		if html[idx] == '"' {
			break
		}
		buf.WriteByte(html[idx])
	}

	data, _ := gg.DoReq(http.MethodGet, "https:"+buf.String(), nil, nil)
	str := string(data)
	befIdx := strings.Index(str, `<script type="text/javascript">`)
	for befIdx > -1 {
		aftIdx := strings.Index(str[befIdx:], "</script>")
		if !strings.Contains(str[befIdx:befIdx+aftIdx], "geturl") {
			str = str[befIdx+aftIdx:]
			befIdx = strings.Index(str, `<script type="text/javascript">`)
		} else {
			vm := goja.New()
			str = strings.ReplaceAll(str[befIdx:befIdx+aftIdx], `<script type="text/javascript">`, "")
			str = strings.ReplaceAll(str, "</script>", "")
			arr := strings.Split(str, "\r\n")
			buf := bytes.Buffer{}
			for _, a := range arr {
				if !strings.Contains(a, "$(") {
					buf.WriteString(a)
					buf.WriteString("\r\n")
				}
			}
			_, err := vm.RunString(buf.String())
			if err != nil {
				return "", err
			}
			v, err := vm.RunString("geturl(gid, currentp)")
			if err != nil {
				return "", err
			}
			return "https:" + v.String(), nil
		}
	}
	return "", nil
}

func (gg *GGBases) GetItemIssueeDate(node *goquery.Document) (time.Time, error) {
	v := node.Find("#udate").AttrOr("v", "")
	unix, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(unix/1000, 0), nil
}

func (gg *GGBases) GetItemLink(node *goquery.Document, id string) ([]model.Link, error) {
	links := []model.Link{}
	node.Find("#showCoverBtn").NextAll().Each(func(i int, s *goquery.Selection) {
		if s.Is("a") {
			url := s.AttrOr("href", "")
			if url != "" {
				if url[:6] != "https:" {
					url = "https:" + url
				}
				links = append(links, model.Link{
					Name: s.Text(),
					Url:  url,
				})
			}
		} else {
			s.Find("a").Each(func(i int, s *goquery.Selection) {
				url := s.AttrOr("href", "")
				if url != "" {
					if url[:6] != "https:" {
						url = "https:" + url
					}
					links = append(links, model.Link{
						Name: s.Text(),
						Url:  url,
					})
				}
			})
		}

	})

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("id", "120270")
	err := writer.Close()
	if err != nil {
		return links, err
	}

	body, err := GGBasesScraper.DoReq(http.MethodPost, GGBasesDetailUri, map[string]string{
		"Content-Type": writer.FormDataContentType(),
		"Referer":      fmt.Sprintf("https://ggbases.dlgal.com/view.so?id=%s", id),
	}, payload)

	if err != nil {
		return links, err
	}
	resp := new(DetailResp)
	err = json.Unmarshal(body, resp)
	if err != nil {
		return links, err
	}
	if resp.HomeUrl != "" {
		links = append(links, model.Link{
			Name: "home",
			Url:  resp.HomeUrl,
		})
	}

	return links, nil
}

func (gg *GGBases) GetItemOtherInfo(node *goquery.Document) (string, error) {
	return node.Find("#description div[markdown-text]").Html()
}
