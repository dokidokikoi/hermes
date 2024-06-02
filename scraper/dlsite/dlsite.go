package dlsite

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
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	comm_tools "github.com/dokidokikoi/go-common/tools"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

var (
	DlSiteDomain    = "https://www.dlsite.com/"
	DlSiteSearchUrl = []string{"https://www.dlsite.com/pro/fsr/=/language/jp/sex_category%5B0%5D/male/keyword/", "keyword", "/work_category%5B0%5D/pc/order%5B0%5D/trend/options_and_or/and/options%5B0%5D/JPN/options%5B1%5D/ENG/options%5B2%5D/CHI/options%5B3%5D/OTL/options%5B4%5D/NM/options_name%5B0%5D/%E6%97%A5%E8%AF%AD%E4%BD%9C%E5%93%81/options_name%5B1%5D/%E8%8B%B1%E8%AF%AD%E4%BD%9C%E5%93%81/options_name%5B2%5D/%E4%B8%AD%E6%96%87%E4%BD%9C%E5%93%81/options_name%5B3%5D/%E4%B8%8D%E9%99%90%E8%AF%AD%E8%A8%80/per_page/30/page/", "page", "/show_type/3/lang_options%5B0%5D/%E6%97%A5%E6%96%87/lang_options%5B1%5D/%E8%8B%B1%E6%96%87/lang_options%5B2%5D/%E4%B8%AD%E6%96%87/lang_options%5B3%5D/%E5%85%B6%E4%BB%96%E8%AF%AD%E8%A8%80/lang_options%5B4%5D/%E4%B8%8D%E9%99%90%E8%AF%AD%E7%A7%8D"}
	DlSitePriceUrl  = "https://www.dlsite.com/pro/product/info/ajax?product_id=%s&cdn_cache_min=1"
	DlSiteVedioApi  = "https://chobit.cc/api/v1/dlsite/embed?workno=%s&_=%d"
)

type DlSite struct {
	Name      string
	Domain    string
	SearchUri string
	Headers   map[string]string
}

var DlSiteScraper *DlSite

func (ds *DlSite) Sreach(keyword string, page int) ([]*scraper.SearchItem, error) {
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
		if len(url) > 6 && url[:6] != "https:" {
			url = "https:" + url
		}
		items = append(items, &scraper.SearchItem{
			Name:  s.Find("dd.work_name a").Text(),
			URl:   s.Find("dl dt a").First().AttrOr("href", ""),
			Cover: url,

			ScraperName: ds.Name,
		})
	})

	return items, nil
}

func (ds *DlSite) DoReq(method, uri string, header map[string]string, body interface{}) ([]byte, error) {
	h := map[string]string{}
	for k, v := range ds.Headers {
		h[k] = v
	}
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
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		r = bytes.NewBuffer(data)
	}

	data, _, err := tools.MakeRequest(method, uri, config.GetConfig().ProxyConfig, r, h, nil)
	return data, err
}

func (ds *DlSite) GetItem(uri string) (*scraper.GameItem, error) {
	data, err := ds.DoReq(http.MethodGet, uri, nil, nil)
	if err != nil {
		return nil, err
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	item := &scraper.GameItem{GameVo: handler.GameVo{Links: []model.Link{{Name: "dlsite", Url: uri}}}}

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
		zaplog.L().Error("获取封面失败", zap.String("scraper", ds.Name), zap.String("uri", uri), zap.Error(err))
	}
	item.Publisher, err = ds.GetItemPublisher(root)
	if err != nil {
		zaplog.L().Error("获取商标失败", zap.String("scraper", ds.Name), zap.String("uri", uri), zap.Error(err))
	}
	item.Characters, err = ds.GetItemCharacter(root)
	if err != nil {
		zaplog.L().Error("获取角色失败", zap.String("scraper", ds.Name), zap.String("uri", uri), zap.Error(err))
	}
	item.Story, err = ds.GetItemStory(root)
	if err != nil {
		zaplog.L().Error("获取故事失败", zap.String("scraper", ds.Name), zap.String("uri", uri), zap.Error(err))
	}
	links, err := ds.GetItemlink(root, id)
	if err != nil {
		zaplog.L().Error("获取链接失败", zap.String("scraper", ds.Name), zap.String("uri", uri), zap.Error(err))
	}
	item.Links = append(item.Links, links...)

	nameSet := map[string]int{}
	root.Find("#work_outline tr").Each(func(i int, s *goquery.Selection) {
		label := s.Find("th").Text()
		if strings.Contains(label, "贩卖日") {
			item.IssueDate = tools.Str2Time(s.Find("td").First().Text())
		} else if strings.Contains(label, "剧情") {
			s.Find("td").First().Each(func(i int, s *goquery.Selection) {
				name := strings.TrimSpace(s.Text())
				if idx, ok := nameSet[name]; ok {
					item.Staff[idx].Relation = append(item.Staff[idx].Relation, model.PRelationWriter.String())
				} else {
					item.Staff = append(item.Staff, handler.StaffVo{
						Name:     name,
						Relation: []string{model.PRelationWriter.String()},
					})
				}
			})
		} else if strings.Contains(label, "插画") {
			s.Find("td").First().Find("a").Each(func(i int, s *goquery.Selection) {
				name := strings.TrimSpace(s.Text())
				if idx, ok := nameSet[name]; ok {
					item.Staff[idx].Relation = append(item.Staff[idx].Relation, model.PRelationPainter.String())
				} else {
					item.Staff = append(item.Staff, handler.StaffVo{
						Name:     name,
						Relation: []string{model.PRelationPainter.String()},
					})
				}
			})
		} else if strings.Contains(label, "声优") {
			s.Find("td").First().Find("a").Each(func(i int, s *goquery.Selection) {
				name := strings.TrimSpace(s.Text())
				if idx, ok := nameSet[name]; ok {
					item.Staff[idx].Relation = append(item.Staff[idx].Relation, model.PRelationCV.String())
				} else {
					item.Staff = append(item.Staff, handler.StaffVo{
						Name:     name,
						Relation: []string{model.PRelationCV.String()},
					})
				}
			})
		} else if strings.Contains(label, "音乐") {
			s.Find("td").First().Find("a").Each(func(i int, s *goquery.Selection) {
				name := strings.TrimSpace(s.Text())
				if idx, ok := nameSet[name]; ok {
					item.Staff[idx].Relation = append(item.Staff[idx].Relation, model.PRelationMusic.String())
				} else {
					item.Staff = append(item.Staff, handler.StaffVo{
						Name:     name,
						Relation: []string{model.PRelationMusic.String()},
					})
				}
			})
		} else {
			s.Find("td").First().Find("a").Each(func(i int, s *goquery.Selection) {
				item.Tags = append(item.Tags, model.Tag{
					Name: strings.TrimSpace(s.Text()),
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
		if len(url) > 6 && url[:6] != "https:" {
			url = "https:" + url
		}
		images = append(images, url)
	})
	if len(images) < 1 {
		return "", nil, nil
	}
	return images[0], images[1:], nil
}

func (ds *DlSite) GetItemPublisher(node *goquery.Document) (*model.Publisher, error) {
	return &model.Publisher{
		Name: node.Find("#work_maker span.maker_name").Text(),
	}, nil
}

func (ds *DlSite) GetItemStory(node *goquery.Document) (string, error) {
	story := bytes.Buffer{}
	node.Find("div.work_parts_container div.work_parts.type_text").Each(func(i int, s *goquery.Selection) {
		html, _ := s.Html()
		story.WriteString(html)
		story.WriteByte('\n')
	})

	return story.String(), nil
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
		if len(url) > 6 && url[:6] != "https:" {
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
		cName := strings.TrimSpace(name.String())

		idx = strings.Index(text, "（CV")
		name.Reset()
		if idx > -1 {
			idx += 5
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
			Summary: text,
			CV: handler.StaffVo{
				Name: strings.TrimSpace(name.String()),
			},
		})
	})

	return characters, nil
}

func init() {
	headers := make(map[string]string)
	headers["Sec-Ch-Ua"] = `"Google Chrome";v="125", "Chromium";v="125", "Not.A/Brand";v="24"`
	headers["Sec-Ch-Ua-Mobile"] = "?0"
	headers["Sec-Ch-Ua-Platform"] = `"macOS"`
	headers["User-Agent"] = config.DefaultUserAgent
	headers["Referer"] = "https://www.getchu.com/php/search.phtml?search_keyword=%C8%E0%BD%F7&list_count=30&sort=sales&sort2=down&search_title=&search_brand=&search_person=&search_jan=&search_isbn=&genre=pc_soft&start_date=&end_date=&age=&list_type=list&search=1&pageID=1"
	headers["Accept-Language"] = config.ZhLanguage
	headers["Cookie"] = "__DLsite_SID=782pmg62psm037ve711d3mcvou; _vwo_uuid_v2=D58EED3043B8C8712836CA3A0CEE347EA|60e7665c582ba675b0d38e5d2fff3d4a; _gcl_au=1.1.1614214747.1716352480; uniqid=0.1jznx3ayl8r; _inflow_ad_params=%7B%22ad_name%22%3A%22organic%22%7D; _fbp=fb.1.1716352481206.123734068; _gaid=876588495.1716352481; _yjsu_yjad=1716352481.85b8a4d8-d0ef-40c3-9add-72cdbb9aefb2; __lt__cid=d543a8ce-a077-4fc3-9c2e-8f3174d95c59; localesuggested=true; locale=zh-cn; _tt_enable_cookie=1; _ttp=rjrtRsh6ouv0PmOqGCRhvWVWaIs; _im_vid=01HYF98ZMNGPY3Z39Q5YJMZW7W; universe_aid=bcf505d16a92b2a620515be740e116240a1a00eccd6e9b0e; adr_id=S7YehhRFnRk3O6gLUtpysCXLJ0EzAGzX1yWf6W6kB4FhT3yt; adultchecked=1; _inflow_params=%7B%22referrer_uri%22%3A%22www.google.com.hk%22%7D; _gid=GA1.2.771398923.1717239817; _ga_QEETZHFB1S=GS1.1.1717290605.1.1.1717290605.0.0.0; _ga_YG879NVEC7=GS1.1.1717290602.1.1.1717290636.0.0.0; _ga_sid=1717304107; __lt__sid=a4160a04-dddad713; DL_PRODUCT_LOG=%2CVJ011538%2CVJ01001190%2CVJ01001393%2CVJ01002056; OptanonConsent=isGpcEnabled=0&datestamp=Sun+Jun+02+2024+14%3A08%3A57+GMT%2B0800+(%E4%B8%AD%E5%9B%BD%E6%A0%87%E5%87%86%E6%97%B6%E9%97%B4)&version=6.23.0&isIABGlobal=false&hosts=&landingPath=NotLandingPage&groups=C0001%3A1%2CC0002%3A0%2CC0003%3A0%2CC0004%3A0&AwaitingReconsent=false; _ga_ZW5GTXK6EV=GS1.1.1717304107.5.1.1717308538.0.0.0; _ga=GA1.1.876588495.1716352481; _inflow_dlsite_params=%7B%22dlsite_referrer_url%22%3A%22https%3A%2F%2Fwww.dlsite.com%2Fpro%2Fwork%2F%3D%2Fproduct_id%2FVJ01001190.html%22%7D"
	DlSiteScraper = &DlSite{
		Name:      "getchu",
		Domain:    DlSiteDomain,
		SearchUri: "",
		Headers:   headers,
	}
}
