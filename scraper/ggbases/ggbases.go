package ggbases

import (
	"bytes"
	"encoding/json"
	"fmt"
	"izumi/internal/handler"
	"izumi/model"
	"izumi/scraper"
	"izumi/tools"
	"maps"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dokidokikoi/go-common/gopool"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	comm_tools "github.com/dokidokikoi/go-common/tools"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const Name = "ggbases"

const (
	DefaultHeader_Cookie = "GGBASESSESSIONID=83D39E9973F44502AFAE0E648A2EF263"
)

var GGBasesDomain = "https://ggbases.dlgal.com/"
var GGBasesSearchUri = "https://ggbases.dlgal.com/search.so?p=%d&title=%s&advanced=0"
var GGBasesMagnetUri = "https://ggb.dlgal.com/magnet.so?id=%s"
var GGBasesBtUri = "https://ggbases.dlgal.com/down.so?id=%s"
var GGBasesDetailUri = "https://ggbases.dlgal.com/viewDetail.so"

type GGBases struct {
	sync.RWMutex
	name      string
	Domain    string
	SearchUri string
	Headers   map[string]string
	Proxy     string
	logger    *zap.Logger

	getitemChan chan struct{}
}

func NewGGBases(header map[string]string, proxy string) scraper.IGameScraper {
	if header == nil {
		header = make(map[string]string)
	}
	header["Referer"] = GGBasesDomain
	return &GGBases{
		name:        Name,
		Domain:      GGBasesDomain,
		SearchUri:   GGBasesSearchUri,
		Headers:     header,
		Proxy:       proxy,
		logger:      zaplog.L().Named(Name),
		getitemChan: make(chan struct{}, 1),
	}
}

func (gg *GGBases) GetName() string {
	return gg.name
}

func (gg *GGBases) SetHeader(header map[string]string) {
	gg.Lock()
	maps.Copy(gg.Headers, header)
	gg.Unlock()
}

func (gg *GGBases) SetProxy(proxy string) {
	gg.Lock()
	gg.Proxy = proxy
	gg.Unlock()
}

func (gg *GGBases) Search(keyword string, page int) ([]*scraper.SearchItem, error) {
	data, err := gg.DoReq(http.MethodGet, fmt.Sprintf(gg.SearchUri, page-1, keyword), map[string]string{
		"Accept-Language": "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7",
	}, nil)
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
				Name:        comm_tools.TrimBlankChar(s.Find("td.t-l").Eq(1).Text()),
				ScraperName: gg.name,
			}

			item.Name = s.Find("a[name='title']").Text()

			item.Cover, err = gg.GetListItemCover(s, item.URl)
			if err != nil {
				gg.logger.Error("获取封面失败", zap.String("scraper", gg.name), zap.Error(err))
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
	maps.Copy(h, gg.Headers)
	gg.RUnlock()
	maps.Copy(h, header)

	rsp, err := tools.ReqWithProxy(method, uri, body, gg.Proxy, tools.SetHeadersWithOption(h))
	if err != nil {
		return nil, err
	}
	return rsp.Bytes(), nil
}

func (gg *GGBases) GetItem(uri string) (*scraper.GameItem, error) {
	// 防止同时请求
	gg.getitemChan <- struct{}{}
	go func() {
		time.Sleep(time.Second * 2)
		<-gg.getitemChan
	}()
	data, err := gg.DoReq(http.MethodGet, uri, nil, nil)
	if err != nil {
		return nil, err
	}
	item := &scraper.GameItem{GameVo: handler.GameVo{Links: []model.Link{{Name: "ggbases", Url: uri}}}, ScraperName: gg.name}
	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, errors.Wrapf(err, "uri error")
	}
	id := u.Query().Get("id")
	if id == "" {
		return nil, errors.New("uri error")
	}

	item.Name, err = gg.GetItemName(root)
	if err != nil {
		gg.logger.Warn("获取名称失败", zap.String("scraper", gg.name), zap.String("uri", uri), zap.Error(err))
	}
	item.Cover, err = gg.GetItemCover(root.Selection)
	if err != nil {
		gg.logger.Warn("获取封面失败", zap.String("scraper", gg.name), zap.String("uri", uri), zap.Error(err))
	}
	item.IssueDate, err = gg.GetItemIssueeDate(root)
	if err != nil {
		gg.logger.Warn("获取发布时间失败", zap.String("scraper", gg.name), zap.String("uri", uri), zap.Error(err))
	}
	links, err := gg.GetItemLink(root, id)
	if err != nil {
		gg.logger.Warn("获取相关链接失败", zap.String("scraper", gg.name), zap.String("uri", uri), zap.Error(err))
	}
	item.Links = append(item.Links, links...)
	item.OtherInfo, err = gg.GetItemOtherInfo(root)
	if err != nil {
		gg.logger.Warn("获取其它信息失败", zap.String("scraper", gg.name), zap.String("uri", uri), zap.Error(err))
	}
	magnet, err := gg.GetMagnet(id)
	if err != nil {
		gg.logger.Warn("获取磁力失败", zap.String("scraper", gg.name), zap.String("uri", uri), zap.Error(err))
	} else {
		item.DownloadInfos = append(item.DownloadInfos, model.DownloadInfo{
			Content: magnet,
			Type:    model.DownloadMagnet,
		})
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
							Relation: []model.PersonRelation{model.PRelationPainter},
						})
					case "group":
						item.Brand = &model.Brand{
							Name: content,
						}
					default:
						item.Tags = append(item.Tags, &model.Tag{
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

func (gg *GGBases) GetListItemCover(node *goquery.Selection, detailUrl string) (string, error) {
	c := node.Find("a[name='title']").AttrOr("c", "")
	cover := GetCover(c)
	// if cover == "" {
	// 	itemData, err := gg.DoReq(http.MethodGet, detailUrl, nil, nil)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(itemData))
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	cover, err = gg.GetItemCover(root.Selection)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// }
	return cover, nil
}

func (gg *GGBases) GetItemCover(root *goquery.Selection) (string, error) {
	cover := root.Find("#ecover").Eq(0).AttrOr("src", "")
	if cover == "" {
		html, err := root.Html()
		if err != nil {
			return "", err
		}
		idx := strings.Index(html, `$("#showCoverBtn").attr("href", "//www.galpic.xyz/show/`)
		sb := strings.Builder{}
		if idx > -1 {
			idx += len(`$("#showCoverBtn").attr("href", "//www.galpic.xyz/show/`)
			for {
				if html[idx] == '"' {
					break
				} else {
					sb.WriteByte(html[idx])
				}
				idx++
			}
		}
		did := sb.String()
		cover = GetCover(did)
		if cover == "" {
			images := MarkdownImg(root.Find("div[markdown-text]").Text())
			if len(images) > 0 {
				cover = images[0]
			}
		}
	}

	return cover, nil
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
				if !strings.HasPrefix(url, "https:") {
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
					if !strings.HasPrefix(url, "https:") {
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

	body, err := gg.DoReq(http.MethodPost, GGBasesDetailUri, map[string]string{
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

func (gg *GGBases) GetMagnet(id string) (string, error) {
	data, err := gg.DoReq(http.MethodGet, fmt.Sprintf(GGBasesMagnetUri, id), map[string]string{
		"Accept":  "application/json",
		"Origin":  "https://ggb.dlgal.com",
		"referer": fmt.Sprintf(GGBasesMagnetUri, id),
	}, nil)
	if err != nil {
		return "", err
	}
	idx := bytes.Index(data, []byte("/magnet.so?"))
	if idx < 0 {
		return "", err
	}
	magnetUrl := ""
	for i := idx; i < len(data); i++ {
		if data[i] == '"' {
			magnetUrl = "https://ggb.dlgal.com" + string(data[idx:i])
			break
		}
	}

	body, err := gg.DoReq(http.MethodPost, magnetUrl, map[string]string{
		"Origin":  "https://ggb.dlgal.com",
		"Referer": fmt.Sprintf(GGBasesMagnetUri, id),
	}, nil)
	if err != nil {
		return "", err
	}

	result := new(struct {
		Hash string `json:"hash"`
	})
	err = json.Unmarshal(body, result)
	if err != nil {
		return "", err
	}
	return "magnet:?xt=urn:btih:" + result.Hash, nil
}
