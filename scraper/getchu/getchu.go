package getchu

import (
	"bytes"
	"fmt"
	"izumi/internal/handler"
	"izumi/model"
	"izumi/scraper"
	"izumi/tools"
	"maps"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/PuerkitoBio/goquery"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	comm_tools "github.com/dokidokikoi/go-common/tools"
	"go.uber.org/zap"
)

const Name = "getchu"

const (
	DefaultHeader_SecChUa         = `"Google Chrome";v="125", "Chromium";v="125", "Not.A/Brand";v="24"`
	DefaultHeader_SecChUaPlatform = `"macOS"`
	DefaultHeader_SecChUaMobile   = "?0"
	DefaultHeader_Referer         = "https://www.getchu.com/php/search.phtml?search_keyword=%C8%E0%BD%F7&list_count=30&sort=sales&sort2=down&search_title=&search_brand=&search_person=&search_jan=&search_isbn=&genre=pc_soft&start_date=&end_date=&age=&list_type=list&search=1&pageID=1"
	DefaultHeader_Cookie          = "_im_vid=01HYF9KCRA1MT8HSM4EWETGX8S; _gid=GA1.2.1781699859.1717215574; getchu_adalt_flag=getchu.com; ITEM_HISTORY=1282568%7C1273918; _ga_BSNR8334HV=GS1.1.1717222828.5.1.1717225315.53.0.0; _ga_JBMY6G3QFS=GS1.1.1717222828.5.1.1717225315.53.0.0; _ga=GA1.2.1343565952.1716352800; _gat=1"
)

var (
	GetChuDomain = "https://www.getchu.com/"
	GetChuSearch = "https://www.getchu.com/php/search.phtml?search_keyword=&list_count=30&sort=sales&sort2=down&search_title=%s&search_brand=&search_person=&search_jan=&search_isbn=&genre=pc_soft&start_date=&end_date=&age=&list_type=list&search=search&pageID=%d"
	GetChuImage  = "https://www.getchu.com/soft_sampleimage.phtml?id=%s"
)

type GetChu struct {
	sync.RWMutex
	name      string
	Domain    string
	SearchUri string
	Headers   map[string]string
	Proxy     string
	logger    *zap.Logger
}

func NewGetChu(headers map[string]string, proxy string) scraper.IGameScraper {
	return &GetChu{
		name:      "getchu",
		Domain:    GetChuDomain,
		SearchUri: "",
		Headers:   headers,
		Proxy:     proxy,
		logger:    zaplog.L().Named(Name),
	}
}

func (gc *GetChu) GetName() string {
	return gc.name
}

func (gc *GetChu) SetHeader(header map[string]string) {
	gc.RLock()
	maps.Copy(gc.Headers, header)
	gc.RUnlock()
}

func (gc *GetChu) SetProxy(proxy string) {
	gc.RLock()
	gc.Proxy = proxy
	gc.RUnlock()
}

func (gc *GetChu) Search(keyword string, page int) ([]*scraper.SearchItem, error) {
	keyword, err := tools.Utf82Jp([]byte(keyword))
	if err != nil {
		return nil, err
	}
	keyword = tools.UrlEnc(keyword)
	url := fmt.Sprintf(GetChuSearch, keyword, page)

	data, err := gc.DoReq(http.MethodGet, url, nil, nil)
	if err != nil {
		return nil, err
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	var urls []string
	var items []*scraper.SearchItem
	root.Find("div.search_container ul.display li").Each(func(i int, s *goquery.Selection) {
		s = s.Find("div.content_block")
		name, err := tools.Jp2Utf8([]byte(s.Find("#detail_block td").Eq(0).Find("a").Eq(0).Text()))
		if err != nil {
			gc.logger.Error("jp encode err", zap.Error(err))
		}
		url := s.Find("#package_block a img").AttrOr("data-original", "")
		urls = append(urls, url)
		items = append(items, &scraper.SearchItem{
			Name:        name,
			URl:         gc.AbsUrl(s.Find("#package_block a").AttrOr("href", "")),
			Cover:       url,
			ScraperName: gc.name,
		})
	})

	m := tools.SaveBunchTmpFile(func(url string) ([]byte, error) {
		return gc.DoReq(http.MethodGet, url, nil, nil)
	}, urls)
	for _, item := range items {
		item.Cover = m[item.Cover]
	}

	return items, nil
}

func (gc *GetChu) DoReq(method, uri string, header map[string]string, body interface{}) ([]byte, error) {
	h := map[string]string{}
	gc.Lock()
	maps.Copy(h, gc.Headers)
	gc.Unlock()
	maps.Copy(h, header)

	query := comm_tools.GenQueryParams(body)
	if query != "" {
		uri += "?" + query
	}

	rsp, err := tools.ReqWithProxy(method, uri, body, gc.Proxy, tools.SetHeadersWithOption(h))
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode()/100 != 2 {
		return nil, errors.Errorf("status code: %d, body: %s", rsp.StatusCode(), rsp.String())
	}
	return rsp.Bytes(), nil
}

func (gc *GetChu) AbsUrl(uri string) string {
	return comm_tools.AbsUrl(gc.Domain, uri)
}

func (gc *GetChu) GetItem(uri string) (*scraper.GameItem, error) {
	data, err := gc.DoReq(http.MethodGet, uri, nil, nil)
	if err != nil {
		return nil, err
	}
	item := &scraper.GameItem{GameVo: handler.GameVo{Links: []model.Link{{Name: "getchu", Url: uri}}}, ScraperName: gc.name}
	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, errors.Wrapf(err, "url.Parse(%s)", uri)
	}
	kv := u.Query()
	id := kv.Get("id")
	if id == "" {
		return nil, errors.Errorf("id not found, uri: %s", uri)
	}

	item.Cover, item.Images, err = gc.GetItemCover(root, id)
	if err != nil {
		gc.logger.Error("获取封面失败", zap.String("scraper", gc.name), zap.String("uri", uri), zap.Error(err))
	}
	item.Name, err = gc.GetItemName(root)
	if err != nil {
		gc.logger.Error("获取名称失败", zap.String("scraper", gc.name), zap.String("uri", uri), zap.Error(err))
	}

	root.Find("#soft_table tr:nth-child(2) table tr").Each(func(i int, s *goquery.Selection) {
		title, err := tools.Jp2Utf8([]byte(s.Find("td").First().Text()))
		if err != nil {
			gc.logger.Error("jp 解码错误", zap.Error(err))
			return
		}

		content := s.Find("td").Eq(1).Text()
		content, err = tools.Jp2Utf8([]byte(content))
		if err != nil {
			gc.logger.Error("jp 解码错误", zap.Error(err))
			return
		}

		nameSet := map[string]int{}

		title = comm_tools.TrimBlankChar(title)
		if strings.Contains(title, "ブランド") {
			item.Brand = &model.Brand{
				Name: content,
			}
		} else if strings.Contains(title, "定価") {
			item.Price = content
		} else if strings.Contains(title, "発売日") {
			item.IssueDate = tools.Str2Time(content)
		} else if strings.Contains(title, "JANコード") {
			item.JanCode = content
		} else if strings.Contains(title, "品番") {
			item.Code = content
		} else if strings.Contains(title, "原画") {
			s.Find("td").Eq(1).Find("a").Each(func(i int, s *goquery.Selection) {
				name, err := tools.Jp2Utf8([]byte(s.Text()))
				if err != nil {
					gc.logger.Error("jp 解码错误", zap.Error(err))
					return
				}
				name = comm_tools.TrimBlankChar(name)

				if idx, ok := nameSet[name]; ok {
					item.Staff[idx].Relation = append(item.Staff[idx].Relation, model.PRelationPainter)
				} else {
					item.Staff = append(item.Staff, handler.StaffVo{
						Name:     name,
						Relation: []model.PersonRelation{model.PRelationPainter},
					})
					nameSet[name] = len(item.Staff) - 1
				}
			})
		} else if strings.Contains(title, "シナリオ") {
			s.Find("td").Eq(1).Find("a").Each(func(i int, s *goquery.Selection) {
				name, err := tools.Jp2Utf8([]byte(s.Text()))
				if err != nil {
					gc.logger.Error("jp 解码错误", zap.Error(err))
					return
				}
				name = comm_tools.TrimBlankChar(name)

				if idx, ok := nameSet[name]; ok {
					item.Staff[idx].Relation = append(item.Staff[idx].Relation, model.PRelationWriter)
				} else {
					item.Staff = append(item.Staff, handler.StaffVo{
						Name:     name,
						Relation: []model.PersonRelation{model.PRelationWriter},
					})
					nameSet[name] = len(item.Staff) - 1
				}
			})
		}
	})

	item.Characters, err = gc.GetItemCharacter(root)
	if err != nil {
		gc.logger.Error("获取角色失败", zap.String("scraper", gc.name), zap.String("uri", uri), zap.Error(err))
	}
	item.Story, err = gc.GetItemStory(root)
	if err != nil {
		gc.logger.Error("获取故事失败", zap.String("scraper", gc.name), zap.String("uri", uri), zap.Error(err))
	}

	return item, nil
}

func (gc *GetChu) GetItemName(node *goquery.Document) (string, error) {
	return tools.Jp2Utf8([]byte(comm_tools.TrimBlankChar(node.Find("#soft-title").Text())))
}

func (gc *GetChu) GetItemCover(node *goquery.Document, id string) (cover string, images []string, err error) {
	cover = gc.AbsUrl(node.Find("#soft_table tr").First().Find("a").First().Find("img").AttrOr("src", ""))

	urls := []string{cover}
	defer func() {
		res := tools.SaveBunchTmpFile(func(url string) ([]byte, error) {
			return gc.DoReq(http.MethodGet, url, nil, nil)
		}, urls)

		for _, url := range urls {
			path, ok := res[url]
			if ok {
				images = append(images, path)
			}
		}

		if len(images) > 0 {
			cover = images[0]
			images = images[1:]
		}
	}()

	data, err := gc.DoReq(http.MethodGet, fmt.Sprintf(GetChuImage, id), nil, nil)
	if err != nil {
		return
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		return
	}

	root.Find(".sample_table_cell").Each(func(i int, s *goquery.Selection) {
		image := s.Find("a").AttrOr("href", "")
		if image != "" {
			urls = append(urls, gc.AbsUrl(image))
		}
	})

	return
}

func (gc *GetChu) GetItemStory(node *goquery.Document) (string, error) {
	var story string
	node.Find("div.tabletitle").Each(func(i int, selection *goquery.Selection) {
		title, err := tools.Jp2Utf8([]byte(selection.Text()))
		if err != nil {
			return
		}
		if strings.Contains(title, "ストーリー") {
			story, err = tools.Jp2Utf8([]byte(selection.Next().Text()))
			if err != nil {
				return
			}
			return
		}
	})
	return comm_tools.TrimBlankChar(story), nil
}

func (gc *GetChu) GetItemCharacter(node *goquery.Document) ([]handler.CharacterVo, error) {
	var characters []handler.CharacterVo
	node.Find("div.tabletitle").Each(func(i int, selection *goquery.Selection) {
		title, err := tools.Jp2Utf8([]byte(selection.Text()))
		if err != nil {
			return
		}
		urls := []string{}
		if strings.Contains(title, "キャラクター") {
			trs := selection.Next().Find(`tr`)
			trs.Each(func(i int, selection *goquery.Selection) {
				if selection.Find("hr").Length() > 0 {
					return
				}
				avatar, _ := selection.Find("td:nth-child(1) img").Attr("src")
				name, err := tools.Jp2Utf8([]byte(selection.Find("td:nth-child(2) h2.chara-name").Text()))
				if err != nil {
					gc.logger.With(zap.Error(err)).Error("tools.Jp2Utf8")
				}
				name = comm_tools.TrimBlankChar(name)
				introduction, err := tools.Jp2Utf8([]byte(selection.Find("td:nth-child(2) dd").Text()))
				if err != nil {
					gc.logger.With(zap.Error(err)).Error("tools.Jp2Utf8")
				}
				introduction = comm_tools.TrimBlankChar(introduction)
				image, _ := selection.Find("td:nth-child(3) a").Attr("href")
				characters = append(characters, handler.CharacterVo{
					Name:    name,
					Summary: introduction,
					Cover:   gc.AbsUrl(avatar),
					Images:  []string{gc.AbsUrl(image)},
				})

				urls = append(urls, gc.AbsUrl(image), gc.AbsUrl(avatar))
			})

			res := tools.SaveBunchTmpFile(func(url string) ([]byte, error) {
				return gc.DoReq(http.MethodGet, url, nil, nil)
			}, urls)
			for i := range characters {
				if path, ok := res[characters[i].Cover]; ok {
					characters[i].Cover = path
				}
				var images []string
				for _, image := range characters[i].Images {
					if path, ok := res[image]; ok {
						images = append(images, path)
					}
				}
				characters[i].Images = images
			}
			return
		}
	})

	return characters, nil
}
