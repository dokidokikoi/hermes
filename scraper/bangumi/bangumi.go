package bangumi

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
	"sync"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
	comm_tools "github.com/dokidokikoi/go-common/tools"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

var (
	bangumiToken     = "eWpvGgoL3vjyfO2EvCwjRhfjmIVZsEx4T1VxTl0W"
	bangumiUserAgent = "dokidokikoi/meta-scraper (https://github.com/dokidokikoi/meta-scraper)"

	BangumiDomain        = "https://api.bgm.tv"
	BangumiSearchUri     = BangumiDomain + "/search/subject/%s"
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

const (
	ResponseGroupSmall  = "small"
	ResponseGroupMedium = "medium"
	ResponseGroupLarge  = "large"
)

type Bangumi struct {
	Name      string
	Domain    string
	SearchUri string
	Headers   map[string]string
}

var BangumiScraper *Bangumi

func (b *Bangumi) DoReq(method, uri string, header map[string]string, body interface{}) ([]byte, error) {
	h := map[string]string{}
	for k, v := range b.Headers {
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

func (b *Bangumi) AbsUrl(uri string) string {
	return comm_tools.AbsUrl(b.Domain, uri)
}

func (b *Bangumi) Sreach(keyword string, page int) ([]*scraper.SearchItem, error) {
	param := SearchParam{
		Type:          SubjectTypeGame,
		ResponseGroup: ResponseGroupMedium,
		Start:         (page-1)*defaultPageSize + 1,
		MaxResults:    defaultPageSize,
	}
	data, err := b.DoReq(http.MethodGet, fmt.Sprintf(b.SearchUri, keyword), nil, param)
	if err != nil {
		return nil, err
	}
	var items []*scraper.SearchItem
	for _, i := range gjson.GetBytes(data, "list").Array() {
		items = append(items, &scraper.SearchItem{
			Name:        i.Get("name").String(),
			Key:         strconv.Itoa(int(i.Get("id").Int())),
			URl:         fmt.Sprintf(BangumiSubjectUri, strconv.Itoa(int(i.Get("id").Int()))),
			Summary:     i.Get("summary").String(),
			Cover:       i.Get("images.common").String(),
			ScraperName: b.Name,
		})
	}
	return items, nil
}

func (b *Bangumi) GetItem(uri string) (*scraper.GameItem, error) {
	data, err := b.DoReq(http.MethodGet, uri, nil, nil)
	if err != nil {
		return nil, err
	}

	item := &scraper.GameItem{GameVo: handler.GameVo{Links: []model.Link{{Name: "bangumi", Url: uri}}}}
	// 获取名称
	item.Name = gjson.GetBytes(data, "name").String()
	nameCN := gjson.GetBytes(data, "name_cn").String()
	if nameCN != "" {
		item.Alias = []string{nameCN}
	}
	item.Cover = gjson.GetBytes(data, "images.medium").String()
	item.Story = gjson.GetBytes(data, "summary").String()
	for _, tag := range gjson.GetBytes(data, "tags").Array() {
		item.Tags = append(item.Tags, model.Tag{
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
				Name: "相关网站",
				Url:  info.Get("value").String(),
			})
		}
	}
	id := strconv.Itoa(int(gjson.GetBytes(data, "id").Int()))
	item.Characters, err = b.GetItemCharacters(id)
	if err != nil {
		zaplog.L().Error("get characters error", zap.String("scraper", b.Name), zap.Error(err))
	}
	item.Staff, err = b.GetItemStaff(id)
	if err != nil {
		zaplog.L().Error("get staff error", zap.String("scraper", b.Name), zap.Error(err))
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
				zaplog.L().Error("request error", zap.String("url", fmt.Sprintf(BangumiCharactersInfoUri, strconv.Itoa(int(id)))), zap.Error(err))
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
				Rlation: c.Get("relation").String(),
				Cover:   c.Get("images.medium").String(),
				Gender:  model.GenderMap[cc.Get("gender").String()].String(),
				Summary: cc.Get("summary").String(),
				Alias:   alias,
			}
			arr := c.Get("actors").Array()
			if len(arr) > 0 {
				character.CV = handler.StaffVo{
					Name:    arr[0].Get("name").String(),
					Cover:   arr[0].Get("images.medium").String(),
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

func (b *Bangumi) GetItemStaff(SubjetID string) ([]handler.StaffVo, error) {
	data, err := b.DoReq(http.MethodGet, fmt.Sprintf(BangumiPersonsUri, SubjetID), nil, nil)
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
				zaplog.L().Error("request error", zap.String("url", fmt.Sprintf(BangumiPersonsInfoUri, strconv.Itoa(int(id)))), zap.Error(err))
				return
			}
			ss := gjson.ParseBytes(data)

			relation := model.PRelationMap[s.Get("relation").String()]
			if relation == model.PRelationNone {
				return
			}

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
				Name:     s.Get("name").String(),
				Cover:    s.Get("images.medium").String(),
				Relation: []string{relation.String()},
				Gender:   model.GenderMap[ss.Get("gender").String()].String(),
				Summary:  ss.Get("short_summary").String(),
				Alias:    alias,
			})
			lock.Unlock()
		}()

	}

	return staff, nil
}

func init() {
	headers := make(map[string]string)
	headers["User-Agent"] = bangumiUserAgent
	headers["Authorization"] = fmt.Sprintf("Bearer %s", bangumiToken)
	BangumiScraper = &Bangumi{
		Name:      "bangumi",
		Domain:    BangumiDomain,
		SearchUri: BangumiSearchUri,
		Headers:   headers,
	}
}
