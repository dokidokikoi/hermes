package twodfan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hermes/config"
	"hermes/internal/handler"
	"hermes/model"
	"hermes/scraper"
	"hermes/tools"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	comm_tools "github.com/dokidokikoi/go-common/tools"
	"go.uber.org/zap"
)

var (
	twoDFanDomain    = "https://2dfan.com/"
	twoDFanSearchUri = "https://2dfan.com/subjects%s/search?keyword=%s"
)

type TwoDFan struct {
	Name      string
	Domain    string
	SearchUri string
	Headers   map[string]string
}

var TwoDFanScraper *TwoDFan

func (tdf *TwoDFan) GetName() string {
	return tdf.Name
}

func (tdf *TwoDFan) Sreach(keyword string, page int) ([]*scraper.SearchItem, error) {
	url := ""
	if page > 1 {
		url = fmt.Sprintf(tdf.SearchUri, fmt.Sprintf("/page/%d", page), keyword)
	} else {
		url = fmt.Sprintf(tdf.SearchUri, "", keyword)
	}

	data, err := tdf.DoReq(http.MethodGet, url, nil, nil)
	if err != nil {
		return nil, err
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	var items []*scraper.SearchItem
	root.Find("ul.intro-list li.media").Each(func(i int, s *goquery.Selection) {
		items = append(items, &scraper.SearchItem{
			Name:        s.Find("#content h4.media-heading a").Text(),
			Key:         "",
			URl:         s.Find("h4.media-heading a").AttrOr("href", ""),
			Summary:     "",
			Cover:       s.Find("img.subject-package").AttrOr("src", ""),
			ScraperName: tdf.Name,
		})
	})

	return items, nil
}

func (tdf *TwoDFan) DoReq(method, uri string, header map[string]string, body interface{}) ([]byte, error) {
	h := map[string]string{}
	for k, v := range tdf.Headers {
		h[k] = v
	}
	for k, v := range header {
		h[k] = v
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	data, _, err = tools.MakeRequest(method, uri, config.GetConfig().ProxyConfig, bytes.NewBuffer(data), h, nil)
	return data, err
}

func (tdf *TwoDFan) AbsUrl(uri string) string {
	return comm_tools.AbsUrl(tdf.Domain, uri)
}

func (tdf *TwoDFan) GetItem(uri string) (*scraper.GameItem, error) {
	data, err := tdf.DoReq(http.MethodGet, uri, nil, nil)
	if err != nil {
		return nil, err
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	item := &scraper.GameItem{GameVo: handler.GameVo{Links: []model.Link{{Name: "2dfan", Url: uri}}}}
	// 获取名称
	item.Name, item.Alias, err = tdf.GetItemName(root)
	if err != nil {
		zaplog.L().Error("获取名称失败", zap.String("scraper", tdf.Name), zap.String("uri", uri), zap.Error(err))
	}
	// 获取
	item.Category, err = tdf.GetItemCategory(root)
	if err != nil {
		zaplog.L().Error("获取分类失败", zap.String("scraper", tdf.Name), zap.String("uri", uri), zap.Error(err))
	}
	item.Cover, item.Images, err = tdf.GetItemCover(root)
	if err != nil {
		zaplog.L().Error("获取封面失败", zap.String("scraper", tdf.Name), zap.String("uri", uri), zap.Error(err))
	}
	item.IssueDate, err = tdf.GetItemIssueDate(root)
	if err != nil {
		zaplog.L().Error("获取发布日期失败", zap.String("scraper", tdf.Name), zap.String("uri", uri), zap.Error(err))
	}
	item.Characters, err = tdf.GetItemCharacters(root)
	if err != nil {
		zaplog.L().Error("获取角色失败", zap.String("scraper", tdf.Name), zap.String("uri", uri), zap.Error(err))
	}
	item.Developer, err = tdf.GetItemDeveloper(root)
	if err != nil {
		zaplog.L().Error("获取开发厂商失败", zap.String("scraper", tdf.Name), zap.String("uri", uri), zap.Error(err))
	}
	item.Publisher, err = tdf.GetItemPublisher(root)
	if err != nil {
		zaplog.L().Error("获取发布厂商失败", zap.String("scraper", tdf.Name), zap.String("uri", uri), zap.Error(err))
	}
	item.Tags, err = tdf.GetItemTags(root)
	if err != nil {
		zaplog.L().Error("获取tag失败", zap.String("scraper", tdf.Name), zap.String("uri", uri), zap.Error(err))
	}
	item.OtherInfo, err = tdf.GetItemOtherInfo(root)
	if err != nil {
		zaplog.L().Error("获取其它信息失败", zap.String("scraper", tdf.Name), zap.String("uri", uri), zap.Error(err))
	}
	item.Links, err = tdf.GetItemLinks(root)
	if err != nil {
		zaplog.L().Error("获取链接失败", zap.String("scraper", tdf.Name), zap.String("uri", uri), zap.Error(err))
	}
	item.Story, item.AllImages, err = tdf.GetItemStory(root)
	if err != nil {
		zaplog.L().Error("获取故事失败", zap.String("scraper", tdf.Name), zap.String("uri", uri), zap.Error(err))
	}
	item.Staff, err = tdf.GetItemStaff(root)
	if err != nil {
		zaplog.L().Error("获取staff失败", zap.String("scraper", tdf.Name), zap.String("uri", uri), zap.Error(err))
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
			zaplog.L().Error("fetch iamge error", zap.String("url", url), zap.Error(err))
		} else {
			path, err := tools.SaveTmpFile(filepath.Ext(url), bytes.NewBuffer(data))
			if err != nil {
				zaplog.L().Error("fetch iamge error", zap.String("url", url), zap.Error(err))
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

func (tdf *TwoDFan) GetItemDeveloper(node *goquery.Document) (*model.Developer, error) {
	return nil, nil
}

func (tdf *TwoDFan) GetItemPublisher(node *goquery.Document) (*model.Publisher, error) {
	return &model.Publisher{
		Name: node.Find("#content div.control-group p.tags").First().Find("a").Text(),
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
	root.Find("#content-pagination div.pagination ul li").Each(func(i int, s *goquery.Selection) {
		page, err := strconv.ParseInt(s.Find("a").Text(), 10, 32)
		if err != nil {
			return
		}
		if maxPage < int(page) {
			maxPage = int(page)
		}
	})

	story := strings.Builder{}
	var lock sync.Mutex
	wait := sync.WaitGroup{}
	var images []string
	f := func(root *goquery.Document, referer string) {
		root.Find("#topic-content img").Each(func(i int, s *goquery.Selection) {
			imgUrl := s.AttrOr("src", "")
			if imgUrl == "" {
				return
			}
			wait.Add(1)
			go func() {
				defer wait.Done()

				data, err := tdf.DoReq(http.MethodGet, imgUrl, map[string]string{"Referer": referer}, nil)
				if err != nil {
					zaplog.L().Error("fetch iamge error", zap.String("url", imgUrl), zap.Error(err))
				} else {
					path, err := tools.SaveTmpFile(filepath.Ext(imgUrl), bytes.NewBuffer(data))
					if err != nil {
						zaplog.L().Error("fetch iamge error", zap.String("url", imgUrl), zap.Error(err))
					} else {
						lock.Lock()
						images = append(images, path)
						lock.Unlock()
					}
				}
			}()

		})

		html, _ := root.Find("#topic-content").Html()
		story.WriteString(html)
		story.WriteByte('\n')
	}

	f(root, tdf.AbsUrl(url))
	for i := 2; i <= maxPage; i++ {
		data, err := tdf.DoReq(http.MethodGet, tdf.AbsUrl(url+fmt.Sprintf("/page/%d", i)), nil, nil)
		if err != nil {
			return "", nil, err
		}

		root, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
		if err != nil {
			return "", nil, err
		}

		f(root, tdf.AbsUrl(url+fmt.Sprintf("/page/%d", i)))
	}
	wait.Wait()

	return story.String(), images, nil
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
			t = tools.Str2Time(tools.DateExtarct(s.Text()))
		}
	})
	return t, nil
}

func (tdf *TwoDFan) GetItemTags(node *goquery.Document) ([]model.Tag, error) {
	var tags []model.Tag
	node.Find("#sidebar div.row-fluid").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Find("div.title").Text(), "常用标签") {
			s.Find("div.block-content a").Each(func(i int, s *goquery.Selection) {
				tags = append(tags, model.Tag{
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
					sta.Relation = append(sta.Relation, model.PRelationPainter.String())
				} else {
					sta = handler.StaffVo{
						Name:     s.Text(),
						Relation: []string{model.PRelationPainter.String()},
					}
				}
				staffMap[s.Text()] = sta
			})
		} else if strings.Contains(s.Text(), "剧本") {
			s.Find("a").Each(func(i int, s *goquery.Selection) {
				sta, ok := staffMap[s.Text()]
				if ok {
					sta.Relation = append(sta.Relation, model.PRelationPainter.String())
				} else {
					sta = handler.StaffVo{
						Name:     s.Text(),
						Relation: []string{model.PRelationWriter.String()},
					}
				}
				staffMap[s.Text()] = sta
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

func init() {
	headers := make(map[string]string)
	headers["User-Agent"] = config.DefaultUserAgent
	headers["Referer"] = twoDFanDomain
	headers["Accept-Language"] = config.ZhLanguage
	headers["Cookie"] = "_ga=GA1.1.566177421.1716285606; pop-blocked=true; _project_hgc_session=amhvTGpZYTdmc3VidU4yQUc2cm01aFdTQzhlTk9NdjI2MXFRVWFsUUw3dmRLTXZ4blYwZ2Q4ZUFOOGtkMld2aTg2YWFtSEpzOFJjTkZSejMvaXg5UytTVzYramdaNzNzbFRXYXJ6a1VLNW5RRzU1L29TK3lyWWJaY0wyVWFKUnN2UDQ0K0hPV2ZDTWx0UFVLdE1tajZ6QndtOGRnWkRndFZIM3BkR0FmaUxVWG5PeGtaeEczRXVWTngvd2hQY25MLS1EbzhJc1ZsbFp3VS92dy8wWGIwWG1nPT0%3D--68accc4aae207d572af489e2c4cfa260efdd5f57; _ga_RF77TZ6QMN=GS1.1.1716638766.7.1.1716641708.0.0.0"
	TwoDFanScraper = &TwoDFan{
		Name:      "2dfan",
		Domain:    twoDFanDomain,
		SearchUri: twoDFanSearchUri,
		Headers:   headers,
	}
}
