package bangumi

import (
	"fmt"
	"hermes/internal/handler"
	"hermes/model"
	"hermes/scraper"
	"hermes/tools"
	"maps"
	"net/http"
	"strconv"
	"sync"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
	comm_tools "github.com/dokidokikoi/go-common/tools"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

const Name = "bangumi"

var (
	DefaultHeader_UserAgent = "dokidokikoi/izumi (https://github.com/dokidokikoi/izumi)"
)

var (
	BangumiDomain        = "https://api.bgm.tv"
	BangumiSearchUri     = BangumiDomain + "/v0/search/subjects"
	BangumiSubjectUri    = BangumiDomain + "/v0/subjects/%s"
	BangumiCharactersUri = BangumiDomain + "/v0/subjects/%s/characters"
	BangumiPersonsUri    = BangumiDomain + "/v0/subjects/%s/persons"

	BangumiCharactersInfoUri = BangumiDomain + "/v0/characters/%s"
	BangumiPersonsInfoUri    = BangumiDomain + "/v0/persons/%s"

	defaultPageSize = 20
)

const (
	SubjectTypeBook = iota + 1
	SubjectTypeAnime
	SubjectTypeMusic
	SubjectTypeGame
	SubjectTypeReal = 6
)

type Bangumi struct {
	sync.RWMutex
	name      string
	Domain    string
	SearchUri string
	Headers   map[string]string
	Proxy     string
	logger    *zap.Logger
}

func NewBangumi(header map[string]string, proxy string) scraper.IGameScraper {
	return &Bangumi{
		name:      Name,
		Domain:    BangumiDomain,
		SearchUri: BangumiSearchUri,
		Headers:   header,
		Proxy:     proxy,
		logger:    zaplog.L().Named(Name),
	}
}

func (b *Bangumi) GetName() string {
	return b.name
}

func (b *Bangumi) SetHeader(header map[string]string) {
	b.Lock()
	maps.Copy(b.Headers, header)
	b.Unlock()
}

func (b *Bangumi) SetProxy(proxy string) {
	b.Lock()
	b.Proxy = proxy
	b.Unlock()
}

func (b *Bangumi) AbsUrl(uri string) string {
	return comm_tools.AbsUrl(b.Domain, uri)
}

func (b *Bangumi) Search(keyword string, page int) ([]*scraper.SearchItem, error) {
	param := SearchParam{
		Filter: SearchParamFilter{
			Nsfw: true,
			Type: []int{SubjectTypeGame},
		},
		Keyword: keyword,
		Limit:   defaultPageSize,
		Offset:  (page - 1) * defaultPageSize,
	}
	data, err := b.DoReq(http.MethodPost, b.SearchUri, nil, param)
	if err != nil {
		return nil, err
	}
	var items []*scraper.SearchItem
	for _, i := range gjson.GetBytes(data, "data").Array() {
		items = append(items, &scraper.SearchItem{
			Name:    i.Get("name").String(),
			Key:     strconv.Itoa(int(i.Get("id").Int())),
			URl:     fmt.Sprintf(BangumiSubjectUri, strconv.Itoa(int(i.Get("id").Int()))),
			Summary: i.Get("summary").String(),
			Cover: func() string {
				cover := i.Get("images.large").String()
				if cover == "" {
					cover = i.Get("images.medium").String()
				}
				if cover == "" {
					cover = i.Get("images.common").String()
				}
				return cover
			}(),
			ScraperName: b.name,
		})
	}
	return items, nil
}

func (b *Bangumi) GetItem(uri string) (*scraper.GameItem, error) {
	data, err := b.DoReq(http.MethodGet, uri, nil, nil)
	if err != nil {
		return nil, err
	}

	item := &scraper.GameItem{GameVo: handler.GameVo{Links: []model.Link{{Name: "bangumi", Url: uri}}}, ScraperName: b.name}
	// 获取名称
	item.Name = gjson.GetBytes(data, "name").String()
	nameCN := gjson.GetBytes(data, "name_cn").String()
	if nameCN != "" {
		item.Alias = []string{nameCN}
	}

	cover := gjson.GetBytes(data, "images.large").String()
	if cover == "" {
		cover = gjson.GetBytes(data, "images.medium").String()
	}
	if cover == "" {
		cover = gjson.GetBytes(data, "images.common").String()
	}
	item.Cover = cover
	item.Story = gjson.GetBytes(data, "summary").String()
	for _, tag := range gjson.GetBytes(data, "tags").Array() {
		item.Tags = append(item.Tags, &model.Tag{
			Name: tag.Get("name").String(),
		})
	}
	for _, info := range gjson.GetBytes(data, "infobox").Array() {
		switch info.Get("key").String() {
		case "别名":
			for _, alias := range info.Get("value").Array() {
				item.Alias = append(item.Alias, alias.Get("v").String())
			}
		case "游戏类型":
			item.Category = &model.Category{
				Name: info.Get("value").String(),
			}
		case "发行日期":
			item.IssueDate = tools.Str2Time(info.Get("value").String())
		case "售价":
			item.Price = info.Get("value").String()
		case "website":
			item.Links = append(item.Links, model.Link{
				Name: "官网",
				Url:  info.Get("value").String(),
			})
		}
	}
	id := strconv.Itoa(int(gjson.GetBytes(data, "id").Int()))
	item.Characters, err = b.GetItemCharacters(id)
	if err != nil {
		b.logger.Error("get characters error", zap.String("scraper", b.name), zap.Error(err))
	}
	item.Staff, err = b.GetItemStaff(id)
	if err != nil {
		b.logger.Error("get staff error", zap.String("scraper", b.name), zap.Error(err))
	}

	return item, nil
}

func (b *Bangumi) GetItemCharacters(SubjetID string) ([]handler.CharacterVo, error) {
	data, err := b.DoReq(http.MethodGet, fmt.Sprintf(BangumiCharactersUri, SubjetID), nil, nil)
	if err != nil {
		return nil, err
	}

	var lock sync.Mutex
	wait := sync.WaitGroup{}
	var characters []handler.CharacterVo
	for _, c := range gjson.ParseBytes(data).Array() {
		c := c
		wait.Add(1)
		go func() {
			defer wait.Done()

			id := c.Get("id").Int()
			data, err := b.DoReq(http.MethodGet, fmt.Sprintf(BangumiCharactersInfoUri, strconv.Itoa(int(id))), nil, nil)
			if err != nil {
				b.logger.Error("request error", zap.String("url", fmt.Sprintf(BangumiCharactersInfoUri, strconv.Itoa(int(id)))), zap.Error(err))
				return
			}
			cc := gjson.ParseBytes(data)

			var alias []string
			for _, a := range cc.Get("infobox").Array() {
				if a.Get("key").String() == "别名" {
					for _, ali := range a.Get("value").Array() {
						alias = append(alias, ali.Get("v").String())
					}
				}
			}
			character := handler.CharacterVo{
				Name:    c.Get("name").String(),
				Rlation: model.CharacterRelation(c.Get("relation").String()),
				Cover: func() string {
					cover := c.Get("images.large").String()
					if cover == "" {
						cover = c.Get("images.medium").String()
					}
					if cover == "" {
						cover = c.Get("images.common").String()
					}
					return cover
				}(),
				Gender:  Gender(cc.Get("gender").String()),
				Summary: cc.Get("summary").String(),
				Alias:   alias,
			}
			arr := c.Get("actors").Array()
			if len(arr) > 0 {
				character.CV = handler.StaffVo{
					Name: arr[0].Get("name").String(),
					Cover: func() string {
						cover := arr[0].Get("images.large").String()
						if cover == "" {
							cover = arr[0].Get("images.medium").String()
						}
						if cover == "" {
							cover = arr[0].Get("images.common").String()
						}
						return cover
					}(),
					Summary: arr[0].Get("short_summary").String(),
				}
			}

			lock.Lock()
			characters = append(characters, character)
			lock.Unlock()
		}()

	}
	wait.Wait()

	return characters, nil
}

func (b *Bangumi) GetItemStaff(SubjectID string) ([]handler.StaffVo, error) {
	data, err := b.DoReq(http.MethodGet, fmt.Sprintf(BangumiPersonsUri, SubjectID), nil, nil)
	if err != nil {
		return nil, err
	}

	var lock sync.Mutex
	wait := sync.WaitGroup{}
	var staff []handler.StaffVo
	for _, s := range gjson.ParseBytes(data).Array() {
		s := s
		wait.Add(1)
		go func() {
			id := s.Get("id").Int()
			data, err := b.DoReq(http.MethodGet, fmt.Sprintf(BangumiPersonsInfoUri, strconv.Itoa(int(id))), nil, nil)
			if err != nil {
				b.logger.Error("request error", zap.String("url", fmt.Sprintf(BangumiPersonsInfoUri, strconv.Itoa(int(id)))), zap.Error(err))
				return
			}
			ss := gjson.ParseBytes(data)

			relation := PersonRelation(s.Get("relation").String())

			var alias []string
			for _, a := range ss.Get("infobox").Array() {
				if a.Get("key").String() == "别名" {
					for _, ali := range a.Get("value").Array() {
						alias = append(alias, ali.Get("v").String())
					}
				}
			}
			lock.Lock()
			staff = append(staff, handler.StaffVo{
				Name: s.Get("name").String(),
				Cover: func() string {
					cover := s.Get("images.large").String()
					if cover == "" {
						cover = s.Get("images.medium").String()
					}
					if cover == "" {
						cover = s.Get("images.common").String()
					}
					return cover
				}(),
				Relation: []model.PersonRelation{relation},
				Gender:   Gender(ss.Get("gender").String()),
				Summary:  ss.Get("short_summary").String(),
				Alias:    alias,
			})
			lock.Unlock()
		}()

	}

	return staff, nil
}
