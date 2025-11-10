package vndb

import (
	"encoding/json"
	"izumi/internal/handler"
	"izumi/model"
	"izumi/scraper"
	"izumi/tools"
	"maps"
	"net/http"
	"strings"
	"sync"
	"time"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const Name = "vndb"

var (
	SearchFields = []string{
		"id", "title", "titles.title", "titles.official", "image.url",
	}
	DetailFields = []string{
		"id", "title", "titles.title", "titles.official",
		"released", "languages", "platforms", "image.url", "description",
		"rating", "screenshots.url", "tags.name", "tags.aliases", "length_minutes",
		"developers.id", "developers.name", "developers.original", "developers.description",
		"developers.extlinks.url", "developers.extlinks.label",
		"staff.id", "staff.role", "staff.description", "staff.gender", "staff.name",
		"staff.original", "staff.extlinks.url", "staff.extlinks.label",
		"va.staff.id", "va.staff.description", "va.staff.gender", "va.staff.name",
		"va.staff.original", "va.staff.extlinks.url", "va.staff.extlinks.label",
		"va.character.id", "va.character.name", "va.character.original", "va.character.description",
		"va.character.image.url", "va.character.blood_type", "va.character.aliases",
		"va.character.height", "va.character.weight", "va.character.bust", "va.character.waist",
		"va.character.hips", "va.character.cup", "va.character.age", "va.character.birthday",
		"va.character.gender", "va.character.traits.name", "va.character.traits.description",
		"va.note", "extlinks.url", "extlinks.label",
	}
)

const (
	DefaultHeader_UserAgent = "dokidokikoi/izumi (https://github.com/dokidokikoi/izumi)"

	VNDBSearchUri = "https://api.vndb.org/kana/vn"
)

type VNDB struct {
	sync.RWMutex
	name    string
	Headers map[string]string
	Proxy   string
	logger  *zap.Logger
}

func NewVNDB(header map[string]string, proxy string) scraper.IGameScraper {
	return &VNDB{
		name:    Name,
		Headers: header,
		Proxy:   proxy,
		logger:  zaplog.L().Named(Name),
	}
}

func (v *VNDB) SetHeader(header map[string]string) {
	v.Lock()
	v.Headers = header
	v.Unlock()
}

func (v *VNDB) SetProxy(proxy string) {
	v.Lock()
	v.Proxy = proxy
	v.Unlock()
}

func (v *VNDB) Search(keyword string, page int) ([]*scraper.SearchItem, error) {
	data, err := v.DoReq(http.MethodPost, VNDBSearchUri, nil, map[string]any{
		"filters": []any{
			"search", "=", keyword,
		},
		"fields":  strings.Join(SearchFields, ","),
		"results": 20,
		"page":    page,
	})
	if err != nil {
		return nil, err
	}
	var resp BaseResponse[VN]
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	items := make([]*scraper.SearchItem, 0, len(resp.Result))
	for _, val := range resp.Result {
		items = append(items, &scraper.SearchItem{
			Name: func() string {
				for _, t := range val.Titles {
					if t.Official {
						return t.Title
					}
				}
				return val.Title
			}(),
			Cover:       val.Image.Url,
			URl:         val.ID,
			ScraperName: v.GetName(),
		})
	}
	return items, nil
}

func (v *VNDB) GetItem(uri string) (*scraper.GameItem, error) {
	data, err := v.DoReq(http.MethodPost, VNDBSearchUri, nil, map[string]any{
		"filters": []any{
			"id", "=", uri,
		},
		"fields": strings.Join(DetailFields, ","),
	})
	if err != nil {
		return nil, err
	}
	var resp BaseResponse[VN]
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	if len(resp.Result) == 0 {
		return nil, nil
	}
	res := resp.Result[0]
	item := &scraper.GameItem{
		GameVo: handler.GameVo{
			VNDBID: res.ID,
			Name: func() string {
				for _, t := range res.Titles {
					if t.Official {
						return t.Title
					}
				}
				return res.Title
			}(),
			Alias: func() []string {
				alias := []string{}
				for _, t := range res.Titles {
					if !t.Official {
						alias = append(alias, t.Title)
					} else {
						alias = append(alias, res.Title)
					}
				}
				return alias
			}(),
			Cover: res.Image.Url,
			Images: func() []string {
				var images []string
				for _, s := range res.Screenshots {
					images = append(images, s.Url)
				}
				return images
			}(),
			Brands: func() []*model.Brand {
				var brands []*model.Brand
				for _, dev := range res.Developers {
					links := []model.Link{}
					for _, extLink := range dev.ExtLinks {
						links = append(links, Link2Link(extLink))
					}
					brands = append(brands, &model.Brand{
						VNDBID: dev.ID,
						Name:   dev.Name,
						Links:  links,
					})
				}
				return brands
			}(),
			IssueDate: time.Now(),
			Story:     res.Description,
			Tags: func() []*model.Tag {
				tags := []*model.Tag{}
				for _, tag := range res.Tags {
					tags = append(tags, &model.Tag{
						Name: tag.Name,
					})
				}
				return tags
			}(),
			Staff: func() []handler.StaffVo {
				staffs := []handler.StaffVo{}
				staffM := map[string]int{}
				for _, staff := range res.Staff {
					idx, ok := staffM[staff.ID]
					if !ok {
						staffs = append(staffs, Staff2Staff(staff))
						staffM[staff.ID] = len(staffs) - 1
					} else if r := TransfRelation(staff.Role); r != model.PRelationUnknown {
						staffs[idx].Relation = append(staffs[idx].Relation, r)
					}
				}
				for _, va := range res.VA {
					idx, ok := staffM[va.Staff.ID]
					if !ok {
						staff := Staff2Staff(va.Staff)
						staff.Relation = append(staff.Relation, model.PRelationCV)
						staffs = append(staffs, staff)
						staffM[va.Staff.ID] = len(staffs) - 1
					} else {
						staffs[idx].Relation = append(staffs[idx].Relation, model.PRelationCV)
					}
				}
				return staffs
			}(),
			Characters: func() []handler.CharacterVo {
				characters := []handler.CharacterVo{}
				for _, va := range res.VA {
					char := va.Character
					c := handler.CharacterVo{
						VNDBID: char.ID,
						Name: func() string {
							if char.Original != "" {
								return char.Original
							}
							return char.Name
						}(),
						Alias: func() []string {
							if char.Original != "" {
								return append(char.Aliases, char.Name)
							}
							return char.Aliases
						}(),
						Gender: func() model.Gender {
							if len(char.Gender) > 1 {
								return TransfGender(char.Gender[1])
							} else if len(char.Gender) > 0 {
								return TransfGender(char.Gender[0])
							}
							return model.UnKnown
						}(),
						Summary: char.Description,
						Cover:   char.Image.Url,
						CV:      Staff2Staff(va.Staff),
						PersonalInfo: model.PersonalInfo{
							Age: char.Age,
							Birthday: func() [2]int {
								if len(char.Birthday) > 1 {
									return [2]int{char.Birthday[0], char.Birthday[1]}
								} else if len(char.Birthday) > 0 {
									return [2]int{char.Birthday[0], 0}
								}
								return [2]int{0, 0}
							}(),
							BloodType: char.BloodType,
							Bust:      char.Bust,
							Cup:       char.Cup,
							Height:    char.Height,
							Weight:    char.Weight,
							Waist:     char.Waist,
							Hips:      char.Hips,
							Traits: func() []model.Trait {
								traits := []model.Trait{}
								for _, trait := range char.Traits {
									traits = append(traits, model.Trait{
										Name:        trait.Name,
										Description: trait.Description,
									})
								}
								return traits
							}(),
						},
					}
					characters = append(characters, c)
				}
				return characters
			}(),
			Links: func() []model.Link {
				links := []model.Link{}
				for _, item := range res.ExtLinks {
					links = append(links, Link2Link(item))
				}
				return links
			}(),
		},
		ScraperName: v.GetName(),
	}
	return item, nil
}

func (v *VNDB) DoReq(method, uri string, header map[string]string, body any) ([]byte, error) {
	h := map[string]string{}
	v.RLock()
	maps.Copy(h, v.Headers)
	v.RUnlock()
	maps.Copy(h, header)

	rsp, err := tools.ReqWithProxy(method, uri, body, v.Proxy, tools.SetHeadersWithOption(h))
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode()/100 != 2 {
		return nil, errors.Errorf("status code: %d, body: %s", rsp.StatusCode(), rsp.String())
	}
	return rsp.Bytes(), nil
}

func (v *VNDB) GetName() string {
	return v.name
}
