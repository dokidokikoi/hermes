package dlsite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"izumi/internal/handler"
	"izumi/model"
	"izumi/scraper"
	"izumi/tools"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	comm_tools "github.com/dokidokikoi/go-common/tools"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

const Name = "dlsite"

const (
	DefaultHeader_SecChUa         = `"Google Chrome";v="125", "Chromium";v="125", "Not.A/Brand";v="24"`
	DefaultHeader_SecChUaPlatform = `"macOS"`
	DefaultHeader_SecChUaMobile   = "?0"
	DefaultHeader_Cookie          = "_im_vid=01HYF9KCRA1MT8HSM4EWETGX8S; _gid=GA1.2.1781699859.1717215574; getchu_adalt_flag=getchu.com; ITEM_HISTORY=1282568%7C1273918; _ga_BSNR8334HV=GS1.1.1717222828.5.1.1717225315.53.0.0; _ga_JBMY6G3QFS=GS1.1.1717222828.5.1.1717225315.53.0.0; _ga=GA1.2.1343565952.1716352800; _gat=1"
)

var (
	DlSiteDomain    = "https://www.dlsite.com/"
	DlSiteSearchUrl = []string{"https://www.dlsite.com/maniax/fsr/=/language/jp/sex_category%5B0%5D/male/keyword/", "keyword", "/work_category%5B0%5D/doujin/work_category%5B1%5D/books/work_category%5B2%5D/pc/work_category%5B3%5D/app/order%5B0%5D/trend/options_and_or/and/per_page/30/page/", "page", "/from/fs.header"}
	DlSitePriceUrl  = "https://www.dlsite.com/pro/product/info/ajax?product_id=%s&cdn_cache_min=1"
	DlSiteVedioApi  = "https://chobit.cc/api/v1/dlsite/embed?workno=%s&_=%d"
)

type DlSite struct {
	sync.RWMutex
	name      string
	Domain    string
	SearchUri string
	Headers   map[string]string
	Proxy     string
	logger    *zap.Logger
}

func NewDlSite(header map[string]string, proxy string) scraper.IGameScraper {
	return &DlSite{
		name:      Name,
		Domain:    DlSiteDomain,
		SearchUri: "",
		Headers:   header,
		Proxy:     proxy,
		logger:    zaplog.L().Named(Name),
	}
}

func (ds *DlSite) GetName() string {
	return ds.name
}

func (ds *DlSite) SetHeader(header map[string]string) {
	ds.Lock()
	maps.Copy(ds.Headers, header)
	ds.Unlock()
}

func (ds *DlSite) SetProxy(proxy string) {
	ds.Lock()
	ds.Proxy = proxy
	ds.Unlock()
}

func (ds *DlSite) Search(keyword string, page int) ([]*scraper.SearchItem, error) {
	DlSiteSearchUrl[1] = keyword
	DlSiteSearchUrl[3] = strconv.Itoa(page)
	url := strings.Join(DlSiteSearchUrl, "")

	data, err := ds.DoReq(http.MethodGet, url, nil, nil)
	if err != nil {
		return nil, err
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	var items []*scraper.SearchItem
	root.Find("#search_result_img_box li.search_result_img_box_inner").Each(func(i int, s *goquery.Selection) {
		url := s.Find("dl dt a img").First().AttrOr("src", "")
		if !strings.HasPrefix(url, "https:") {
			url = "https:" + url
		}
		items = append(items, &scraper.SearchItem{
			Name:  s.Find("dd.work_name a").Text(),
			URl:   s.Find("dl dt a").First().AttrOr("href", ""),
			Cover: url,

			ScraperName: ds.name,
		})
	})

	return items, nil
}

func (ds *DlSite) DoReq(method, uri string, header map[string]string, body interface{}) ([]byte, error) {
	h := map[string]string{}
	ds.RLock()
	maps.Copy(h, ds.Headers)
	ds.RUnlock()
	maps.Copy(h, header)

	query := comm_tools.GenQueryParams(body)
	if query != "" {
		uri += "?" + query
	}

	rsp, err := tools.ReqWithProxy(method, uri, body, ds.Proxy, tools.SetHeadersWithOption(h))
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode()/100 != 2 {
		return nil, errors.Errorf("status code: %d, body: %s", rsp.StatusCode(), rsp.String())
	}

	return rsp.Bytes(), nil
}

func (ds *DlSite) GetItem(uri string) (*scraper.GameItem, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Add("locale", "zh_CN")
	u.RawQuery = query.Encode()
	data, err := ds.DoReq(http.MethodGet, u.String(), nil, nil)
	if err != nil {
		return nil, err
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	item := &scraper.GameItem{GameVo: handler.GameVo{Links: []model.Link{{Name: Name, Url: uri}}}, ScraperName: ds.name}

	id := ""
	arr := strings.Split(uri, "/")
	if len(arr) > 2 {
		id = arr[len(arr)-1]
		if len(id) > 5 {
			id = id[:len(id)-5]
		}
	}

	item.Name = root.Find("#work_name").Text()
	item.Code = id
	item.Cover, item.Images, err = ds.GetItemCover(root)
	if err != nil {
		ds.logger.Error("获取封面失败", zap.String("scraper", ds.name), zap.String("uri", uri), zap.Error(err))
	}
	item.Brand, err = ds.GetItemBrand(root)
	if err != nil {
		ds.logger.Error("获取开发商失败", zap.String("scraper", ds.name), zap.String("uri", uri), zap.Error(err))
	}
	item.Characters, err = ds.GetItemCharacter(root)
	if err != nil {
		ds.logger.Error("获取角色失败", zap.String("scraper", ds.name), zap.String("uri", uri), zap.Error(err))
	}
	item.Story, err = ds.GetItemStory(root)
	if err != nil {
		ds.logger.Error("获取故事失败", zap.String("scraper", ds.name), zap.String("uri", uri), zap.Error(err))
	}
	links, err := ds.GetItemlink(root, id)
	if err != nil {
		ds.logger.Error("获取链接失败", zap.String("scraper", ds.name), zap.String("uri", uri), zap.Error(err))
	}
	item.Links = append(item.Links, links...)

	nameSet := map[string]int{}
	root.Find("#work_outline tr").Each(func(i int, s *goquery.Selection) {
		label := s.Find("th").Text()
		if slices.Contains([]string{"贩卖日", "販売日"}, label) {
			item.IssueDate = tools.Str2Time(s.Find("td").First().Text())
		} else if slices.Contains([]string{"剧情", "シナリオ"}, label) {
			s.Find("td").First().Each(func(i int, s *goquery.Selection) {
				name := comm_tools.TrimBlankChar(s.Text())
				if idx, ok := nameSet[name]; ok {
					item.Staff[idx].Relation = append(item.Staff[idx].Relation, model.PRelationWriter)
				} else {
					item.Staff = append(item.Staff, handler.StaffVo{
						Name:     name,
						Relation: []model.PersonRelation{model.PRelationWriter},
					})
				}
			})
		} else if slices.Contains([]string{"插画", "イラスト"}, label) {
			s.Find("td").First().Find("a").Each(func(i int, s *goquery.Selection) {
				name := comm_tools.TrimBlankChar(s.Text())
				if idx, ok := nameSet[name]; ok {
					item.Staff[idx].Relation = append(item.Staff[idx].Relation, model.PRelationPainter)
				} else {
					item.Staff = append(item.Staff, handler.StaffVo{
						Name:     comm_tools.TrimBlankChar(name),
						Relation: []model.PersonRelation{model.PRelationPainter},
					})
				}
			})
		} else if slices.Contains([]string{"声优", "声優"}, label) {
			s.Find("td").First().Find("a").Each(func(i int, s *goquery.Selection) {
				name := comm_tools.TrimBlankChar(s.Text())
				if idx, ok := nameSet[name]; ok {
					item.Staff[idx].Relation = append(item.Staff[idx].Relation, model.PRelationCV)
				} else {
					item.Staff = append(item.Staff, handler.StaffVo{
						Name:     name,
						Relation: []model.PersonRelation{model.PRelationCV},
					})
				}
			})
		} else if slices.Contains([]string{"音乐", "音楽"}, label) {
			s.Find("td").First().Find("a").Each(func(i int, s *goquery.Selection) {
				name := comm_tools.TrimBlankChar(s.Text())
				if idx, ok := nameSet[name]; ok {
					item.Staff[idx].Relation = append(item.Staff[idx].Relation, model.PRelationMusic)
				} else {
					item.Staff = append(item.Staff, handler.StaffVo{
						Name:     name,
						Relation: []model.PersonRelation{model.PRelationMusic},
					})
				}
			})
		} else {
			s.Find("td").First().Find("a").Each(func(i int, s *goquery.Selection) {
				item.Tags = append(item.Tags, &model.Tag{
					Name: comm_tools.TrimBlankChar(s.Text()),
				})
			})
		}
	})

	data, err = ds.DoReq(http.MethodGet, fmt.Sprintf(DlSitePriceUrl, id), nil, nil)
	if err != nil {
		return item, nil
	}
	m := map[string]map[string]any{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return item, nil
	}
	tmp, ok := m[id]
	if ok {
		item.Price, _ = tmp["price_str"].(string)
		category, _ := tmp["work_type"].(string)
		item.Category = &model.Category{
			Name: category,
		}
	}

	return item, nil
}

func (ds *DlSite) GetItemCover(node *goquery.Document) (string, []string, error) {
	images := []string{}
	node.Find("div.product-slider div.product-slider-data div").Each(func(i int, s *goquery.Selection) {
		url := s.AttrOr("data-src", "")
		if !strings.HasPrefix(url, "https:") {
			url = "https:" + url
		}
		images = append(images, url)
	})
	if len(images) < 1 {
		return "", nil, nil
	}
	return images[0], images[1:], nil
}

func (ds *DlSite) GetItemBrand(node *goquery.Document) (*model.Brand, error) {
	return &model.Brand{
		Name: comm_tools.TrimBlankChar(node.Find("#work_maker span.maker_name").Text()),
	}, nil
}

func (ds *DlSite) GetItemStory(node *goquery.Document) (string, error) {
	story := bytes.Buffer{}
	node.Find("div.work_parts_container div.work_parts.type_text").Each(func(i int, s *goquery.Selection) {
		html, _ := s.Html()
		story.WriteString(html)
		story.WriteByte('\n')
	})

	return comm_tools.TrimBlankChar(story.String()), nil
}

func (ds *DlSite) GetItemlink(node *goquery.Document, id string) ([]model.Link, error) {
	data, err := ds.DoReq(http.MethodGet, fmt.Sprintf(DlSiteVedioApi, id, time.Now().UnixMilli()), nil, nil)
	if err != nil {
		return nil, err
	}
	str := string(data)
	str = strings.Replace(str, "response(", "", 1)
	str = str[:len(str)-1]
	links := []model.Link{}
	for _, a := range gjson.Get(str, "works").Array() {
		url := a.Get("embed_url").String()
		if a.Get("file_type").String() == "video" {
			data, err = ds.DoReq(http.MethodGet, a.Get("embed_url").String(), nil, nil)
			if err != nil {
				return links, err
			}
			root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
			if err != nil {
				return nil, err
			}
			url = root.Find("meta[itemprop=contentUrl]").AttrOr("content", "")
		}

		links = append(links, model.Link{
			Name: a.Get("work_name").String(),
			Url:  url,
			Type: model.LinkTypeMap[a.Get("file_type").String()],
		})
	}

	return links, nil
}

func (ds *DlSite) GetItemCharacter(node *goquery.Document) ([]handler.CharacterVo, error) {
	characters := []handler.CharacterVo{}
	node.Find("div.work_parts_container div.work_parts.type_multiimages .work_parts_area .work_parts_multiimage li.work_parts_multiimage_item").Each(func(i int, s *goquery.Selection) {
		url := s.Find(".image a").AttrOr("href", "")
		if !strings.HasPrefix(url, "https:") {
			url = "https:" + url
		}
		text, _ := s.Find(".text").Html()
		idx := strings.Index(text, "●")
		name := bytes.Buffer{}
		if idx > -1 {
			idx += 3
			for ; idx < len(text)-3; idx++ {
				if text[idx:idx+3] == "（" || text[idx:idx+3] == "<br" {
					break
				}
				name.WriteByte(text[idx])
			}
		}
		cName := comm_tools.TrimBlankChar(name.String())

		idx = strings.Index(text, "（CV：")
		name.Reset()
		if idx > -1 {
			idx += 8
			for ; idx < len(text)-3; idx++ {
				if text[idx:idx+3] == "）" {
					break
				}
				name.WriteByte(text[idx])
			}
		}
		characters = append(characters, handler.CharacterVo{
			Name:    cName,
			Cover:   url,
			Summary: comm_tools.TrimBlankChar(text),
			CV: handler.StaffVo{
				Name: comm_tools.TrimBlankChar(name.String()),
			},
		})
	})

	return characters, nil
}
