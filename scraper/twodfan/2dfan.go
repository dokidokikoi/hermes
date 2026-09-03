package twodfan

import (
	"bytes"
	"fmt"
	"izumi/config"
	"izumi/internal/handler"
	"izumi/model"
	"izumi/scraper"
	"izumi/utils"
	"maps"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/dokidokikoi/go-common/tools"
	"go.uber.org/zap"
)

const Name = "2dfan"

const (
	DefaultHeader_Cookie = "_ga=GA1.1.274316804.1751095279; _ga_RF77TZ6QMN=GS2.1.s1751095279$o1$g1$t1751095328$j11$l0$h0; _project_hgc_session=ZnlS6F05146BhfkKsXHOEt3ewMCtzyUjnWY1VZM34Ix31OLzodQDGw2i%2FXIEZJER95WXODpVHKxgiOgruCz2YyVRDqzSPpSVSg7yZIIZiTvKZ3KMaFYeU833GjlKT92JH2%2BrOumlx81gmVZh3eHK3Xl94SwBPeac24%2BfLQKiPgoGJ661KGACNvGX3uAKRhdrJwrsKIaXO4f9w%2BKvtHYOBJHfBTtum10e3nRqTgc%2BtNdfxyDTVb80T%2B0S8G1hrZ%2F%2FyP4PtTUpPLxlb9XS4XBfvDv8jnLMeY94BY%2BiYRcL5MxFQMVqsI0%2FoTBBXJHtaTik7BwV54qFfFCudJ1gAAzT1YtfFxO6SARFusy1aHk%2FWh2zlBegAl3J7Jf3S5M%3D--SFuz5Ln45ktVvec7--hR4VT3NvYGEcup6dWqvRpQ%3D%3D"
)

var (
	twoDFanDomain    = "https://2dfan.com/"
	twoDFanSearchUri = "https://2dfan.com/subjects/search%s?keyword=%s"
)

type TwoDFan struct {
	sync.RWMutex
	name      string
	Domain    string
	SearchUri string
	Headers   map[string]string
	Proxy     string
	logger    *zap.Logger
}

func NewTwoDFan(header map[string]string, proxy string) scraper.IGameScraper {
	return &TwoDFan{
		name:      Name,
		Domain:    twoDFanDomain,
		SearchUri: twoDFanSearchUri,
		Headers:   header,
		Proxy:     proxy,
		logger:    zaplog.L().Named(Name),
	}
}

func (tdf *TwoDFan) GetName() string {
	return tdf.name
}

func (tdf *TwoDFan) SetHeader(header map[string]string) {
	tdf.Lock()
	maps.Copy(tdf.Headers, header)
	tdf.Unlock()
}

func (tdf *TwoDFan) SetProxy(proxy string) {
	tdf.Lock()
	tdf.Proxy = proxy
	tdf.Unlock()
}

func (tdf *TwoDFan) SearchGame(keyword string, page int) ([]*scraper.SearchItem, error) {
	u := ""
	if page > 1 {
		u = fmt.Sprintf(tdf.SearchUri, fmt.Sprintf("/page/%d", page), strings.ReplaceAll(strings.TrimSpace(keyword), " ", "+"))
	} else {
		u = fmt.Sprintf(tdf.SearchUri, "", keyword)
	}

	uri, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	uri.RawQuery = uri.Query().Encode()
	data, err := tdf.DoReq(http.MethodGet, uri.String(), nil, nil)
	if err != nil {
		return nil, err
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	var images []string
	var items []*scraper.SearchItem
	root.Find("ul.intro-list li.media").Each(func(i int, s *goquery.Selection) {
		url := s.Find("img.subject-package").AttrOr("data-normal", "")
		images = append(images, url)
		items = append(items, &scraper.SearchItem{
			Name:        s.Find("#content h4.media-heading a").Text(),
			Key:         "",
			URl:         tdf.AbsUrl(s.Find("h4.media-heading a").AttrOr("href", "")),
			Summary:     "",
			Cover:       url,
			ScraperName: tdf.name,
		})
	})

	m := utils.SaveBunchTmpFile(func(url string) ([]byte, error) {
		return tdf.DoReq(http.MethodGet, url, nil, nil)

	}, images)
	for _, item := range items {
		item.Cover = m[item.Cover]
	}
	return items, nil
}

func (tdf *TwoDFan) DoReq(method, uri string, header map[string]string, body any) ([]byte, error) {
	h := map[string]string{}
	tdf.RLock()
	maps.Copy(h, tdf.Headers)
	tdf.RUnlock()
	maps.Copy(h, header)

	rsp, err := tools.ReqWithProxy(method, uri, body, tdf.Proxy, tools.SetHeadersWithOption(h))
	if err != nil {
		return nil, err
	}
	return rsp.Bytes(), nil
}

func (tdf *TwoDFan) AbsUrl(uri string) string {
	return tools.AbsUrl(tdf.Domain, uri)
}

func (tdf *TwoDFan) GetGameItem(uri string) (*scraper.GameItem, error) {
	data, err := tdf.DoReq(http.MethodGet, uri, nil, nil)
	if err != nil {
		return nil, err
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	item := &scraper.GameItem{GameVo: handler.GameVo{Links: []model.Link{{Name: tdf.name, Url: uri, Type: model.LinkTypeInfo}}}, ScraperName: tdf.name}
	u, err := url.Parse(uri)
	if err == nil {
		arr := strings.Split(u.Path, "/")
		id := arr[len(arr)-1]
		if id != "" {
			item.RelIDs = append(item.RelIDs, fmt.Sprintf("%s:%s", Name, id))
		}
	}
	// 获取名称
	item.Name, item.Alias, err = tdf.GetItemName(root)
	if err != nil {
		tdf.logger.Error("获取名称失败", zap.String("scraper", tdf.name), zap.String("uri", uri), zap.Error(err))
	}
	// 获取
	item.Category, err = tdf.GetItemCategory(root)
	if err != nil {
		tdf.logger.Error("获取分类失败", zap.String("scraper", tdf.name), zap.String("uri", uri), zap.Error(err))
	}
	item.Cover, item.Images, err = tdf.GetItemCover(root)
	if err != nil {
		tdf.logger.Error("获取封面失败", zap.String("scraper", tdf.name), zap.String("uri", uri), zap.Error(err))
	}
	item.IssueDate, err = tdf.GetItemIssueDate(root)
	if err != nil {
		tdf.logger.Error("获取发布日期失败", zap.String("scraper", tdf.name), zap.String("uri", uri), zap.Error(err))
	}
	item.Characters, err = tdf.GetItemCharacters(root)
	if err != nil {
		tdf.logger.Error("获取角色失败", zap.String("scraper", tdf.name), zap.String("uri", uri), zap.Error(err))
	}
	item.Brands, err = tdf.GetItemBrands(root)
	if err != nil {
		tdf.logger.Error("获取开发厂商失败", zap.String("scraper", tdf.name), zap.String("uri", uri), zap.Error(err))
	}
	item.Tags, err = tdf.GetItemTags(root)
	if err != nil {
		tdf.logger.Error("获取tag失败", zap.String("scraper", tdf.name), zap.String("uri", uri), zap.Error(err))
	}
	item.OtherInfo, err = tdf.GetItemOtherInfo(root)
	if err != nil {
		tdf.logger.Error("获取其它信息失败", zap.String("scraper", tdf.name), zap.String("uri", uri), zap.Error(err))
	}
	links, err := tdf.GetItemLinks(root)
	if err != nil {
		tdf.logger.Error("获取链接失败", zap.String("scraper", tdf.name), zap.String("uri", uri), zap.Error(err))
	}
	item.Links = append(item.Links, links...)
	var images []string
	item.Story, images, err = tdf.GetItemStory(root)
	if err != nil {
		tdf.logger.Error("获取故事失败", zap.String("scraper", tdf.name), zap.String("uri", uri), zap.Error(err))
	}
	item.Images = append(item.Images, images...)
	item.Staff, err = tdf.GetItemStaff(root)
	if err != nil {
		tdf.logger.Error("获取staff失败", zap.String("scraper", tdf.name), zap.String("uri", uri), zap.Error(err))
	}

	return item, nil
}

func (tdf *TwoDFan) GetItemName(node *goquery.Document) (string, []string, error) {
	var alias []string
	node.Find("#content div.control-group p.tags").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "又名") {
			alias = append(alias, strings.Split(s.Find("span").Text(), "/")...)
		}
	})
	return node.Find("#content div.navbar h3").First().Text(), alias, nil
}

func (tdf *TwoDFan) GetItemCover(node *goquery.Document) (string, []string, error) {
	url := node.Find("#package-image img.subject-package").First().AttrOr("src", "")
	images := []string{}
	if url != "" {
		data, err := tdf.DoReq(http.MethodGet, url, map[string]string{"Referer": tdf.Domain}, nil)
		if err != nil {
			tdf.logger.Error("fetch iamge error", zap.String("url", url), zap.Error(err))
		} else {
			path, err := utils.SaveTmpFile(filepath.Ext(url), bytes.NewBuffer(data))
			if err != nil {
				tdf.logger.Error("fetch iamge error", zap.String("url", url), zap.Error(err))
			} else {
				images = append(images, path)
			}
		}
	}
	if len(images) < 1 {
		return "", nil, nil
	}
	return images[0], images[1:], nil
}

func (tdf *TwoDFan) GetItemCategory(node *goquery.Document) (*model.Category, error) {
	return nil, nil
}

func (tdf *TwoDFan) GetItemBrands(node *goquery.Document) ([]*model.Brand, error) {
	return []*model.Brand{
		{Name: node.Find("#content div.control-group p.tags").First().Find("a").Text()},
	}, nil
}

func (tdf *TwoDFan) GetItemPrice(node *goquery.Document) (string, error) {
	return "", nil
}

func (tdf *TwoDFan) GetItemStory(node *goquery.Document) (string, []string, error) {
	url := ""
	node.Find("#resources span").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "介绍") {
			url = s.Find("a").AttrOr("href", "")
		}
	})
	if url == "" {
		return "", nil, nil
	}

	data, err := tdf.DoReq(http.MethodGet, tdf.AbsUrl(url), nil, nil)
	if err != nil {
		return "", nil, err
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		return "", nil, err
	}

	maxPage := 0
	href := root.Find("#content-pagination div.pagination ul li").Last().Find("a").AttrOr("href", "")
	hrefArr := strings.Split(href, "/")
	if len(hrefArr) > 2 && hrefArr[len(hrefArr)-2] == "page" {
		maxPage, _ = strconv.Atoi(hrefArr[len(hrefArr)-1])
	}

	story := strings.Builder{}
	var images []string
	f := func(root *goquery.Document) {
		root.Find("#topic-content img").Each(func(i int, s *goquery.Selection) {
			imgUrl := s.AttrOr("src", "")
			if imgUrl == "" {
				return
			}

			images = append(images, imgUrl)
		})

		html, _ := root.Find("#topic-content").Html()
		story.WriteString(html)
		story.WriteByte('\n')
	}

	f(root)
	for i := 2; i <= maxPage; i++ {
		data, err := tdf.DoReq(http.MethodGet, tdf.AbsUrl(url+fmt.Sprintf("/page/%d", i)), nil, nil)
		if err != nil {
			return "", nil, err
		}

		root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
		if err != nil {
			return "", nil, err
		}

		f(root)
	}

	res := utils.SaveBunchTmpFile(func(url string) ([]byte, error) {
		return tdf.DoReq(http.MethodGet, url, nil, nil)
	}, images)
	images = images[:0]
	s := story.String()
	for k, v := range res {
		images = append(images, v)
		s = strings.ReplaceAll(s, k, filepath.Join("/api/file", strings.TrimPrefix(v, config.TmpDir)))
	}
	s = tools.TrimBlankChar(s)

	return s, images, nil
}

func (tdf *TwoDFan) GetItemSize(node *goquery.Document) (string, error) {
	return "", nil
}

func (tdf *TwoDFan) GetItemLinks(node *goquery.Document) ([]model.Link, error) {
	url := ""
	var links []model.Link
	node.Find("#sidebar div.row-fluid").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Find("div.title").Text(), "相关帖子") {
			url = s.Find("div.title small a").AttrOr("href", "")
		}
	})

	if url == "" {
		return nil, nil
	}
	data, err := tdf.DoReq(http.MethodGet, tdf.AbsUrl(url), nil, nil)
	if err != nil {
		return nil, err
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	root.Find("div.block-content table tr").Each(func(i int, s *goquery.Selection) {
		n := s.Find("td").Eq(1).Find("a")
		links = append(links, model.Link{
			Name: n.Text(),
			Url:  tdf.AbsUrl(n.AttrOr("href", "")),
			Type: model.LinkTypeOther,
		})
	})

	return links, nil
}

func (tdf *TwoDFan) GetItemCharacters(node *goquery.Document) ([]handler.CharacterVo, error) {
	return nil, nil
}

func (tdf *TwoDFan) GetItemIssueDate(node *goquery.Document) (time.Time, error) {
	var t time.Time
	node.Find("#content div.control-group p.tags").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "发售日期") {
			t = tools.Str2Time(utils.DateExtarct(s.Text()))
		}
	})
	return t, nil
}

func (tdf *TwoDFan) GetItemTags(node *goquery.Document) ([]*model.Tag, error) {
	var tags []*model.Tag
	node.Find("#sidebar div.row-fluid").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Find("div.title").Text(), "常用标签") {
			s.Find("div.block-content a").Each(func(i int, s *goquery.Selection) {
				tags = append(tags, &model.Tag{
					Name: s.Text(),
				})
			})
		}
	})

	return tags, nil
}

func (tdf *TwoDFan) GetItemStaff(node *goquery.Document) ([]handler.StaffVo, error) {
	var staff []handler.StaffVo
	staffMap := map[string]handler.StaffVo{}
	node.Find("#content div.control-group p.tags").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "原画") {
			s.Find("a").Each(func(i int, s *goquery.Selection) {
				sta, ok := staffMap[s.Text()]
				if ok {
					sta.Relation = append(sta.Relation, model.PRelationPainter)
				} else {
					sta = handler.StaffVo{
						Name:     s.Text(),
						Relation: []model.PersonRelation{model.PRelationPainter},
					}
				}
				staffMap[s.Text()] = sta
			})
		} else if strings.Contains(s.Text(), "剧本") {
			s.Find("a").Each(func(i int, s *goquery.Selection) {
				sta, ok := staffMap[s.Text()]
				if ok {
					sta.Relation = append(sta.Relation, model.PRelationPainter)
				} else {
					sta = handler.StaffVo{
						Name:     s.Text(),
						Relation: []model.PersonRelation{model.PRelationWriter},
					}
				}
				staffMap[s.Text()] = sta
			})
		} else if strings.Contains(s.Text(), "声优") {
			s.Find("a").Each(func(i int, s *goquery.Selection) {
				sta, ok := staffMap[s.Text()]
				if ok {
					sta.Relation = append(sta.Relation, model.PRelationCV)
				} else {
					sta = handler.StaffVo{
						Name:     s.Text(),
						Relation: []model.PersonRelation{model.PRelationCV},
					}
				}
			})
		} else if strings.Contains(s.Text(), "音乐") {
			s.Find("a").Each(func(i int, s *goquery.Selection) {
				sta, ok := staffMap[s.Text()]
				if ok {
					sta.Relation = append(sta.Relation, model.PRelationMusic)
				} else {
					sta = handler.StaffVo{
						Name:     s.Text(),
						Relation: []model.PersonRelation{model.PRelationMusic},
					}
				}
			})
		}
	})

	for _, v := range staffMap {
		staff = append(staff, v)
	}

	return staff, nil
}

func (tdf *TwoDFan) GetItemOtherInfo(node *goquery.Document) (string, error) {
	return "", nil
}
