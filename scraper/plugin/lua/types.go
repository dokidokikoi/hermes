package lua

import (
	"izumi/internal/handler"
	"izumi/model"
	"izumi/scraper"

	"github.com/yuin/gopher-lua"
)

// LuaSearchItem represents a search result from Lua
type LuaSearchItem struct {
	Name        string
	Key         string
	URL         string
	Summary     string
	Cover       string
	ScraperName string
}

// LuaGameItem represents a game item from Lua
type LuaGameItem struct {
	VNDBID     string
	JanCode    string
	DlCode     string
	Name       string
	Alias      []string
	Cover      string
	Images     []string
	Category   *LuaCategory
	Series     []*LuaSeries
	Brands     []*LuaBrand
	Price      string
	IssueDate  string // ISO format
	Story      string
	Platform   string
	Tags       []*LuaTag
	Characters []*LuaCharacter
	Links      []*LuaLink
	OtherInfo  string
	Staff      []*LuaStaff
}

// LuaCategory represents a category from Lua
type LuaCategory struct {
	Name string
}

// LuaSeries represents a series from Lua
type LuaSeries struct {
	Name string
}

// LuaBrand represents a brand from Lua
type LuaBrand struct {
	VNDBID string
	Name   string
	Links  []*LuaLink
}

// LuaTag represents a tag from Lua
type LuaTag struct {
	Name string
	Lang string
}

// LuaCharacter represents a character from Lua
type LuaCharacter struct {
	VNDBID       string
	Name         string
	Alias        []string
	Gender       string
	Relation     string
	Summary      string
	Cover        string
	Images       []string
	CV           *LuaStaff
	Tags         []*LuaTag
	PersonalInfo *LuaPersonalInfo
}

// LuaStaff represents a staff/person from Lua
type LuaStaff struct {
	VNDBID    string
	Name      string
	Alias     []string
	Cover     string
	Images    []string
	Tags      []*LuaTag
	Summary   string
	Gender    string
	Relation  []string
	Links     []*LuaLink
}

// LuaPersonalInfo represents personal info from Lua
type LuaPersonalInfo struct {
	Age       int
	Birthday  []int // [year, month]
	BloodType string
	Bust      int
	Cup       string
	Height    int
	Weight    int
	Waist     int
	Hips      int
	Traits    []*LuaTrait
}

// LuaTrait represents a trait from Lua
type LuaTrait struct {
	Name        string
	Description string
}

// LuaLink represents a link from Lua
type LuaLink struct {
	Name string
	Url  string
	Type string
}

// luaTableToMap converts a Lua table to a map[string]string
func luaTableToMap(L *lua.LState, t *lua.LTable) map[string]string {
	result := make(map[string]string)
	t.ForEach(func(key, value lua.LValue) {
		if strKey, ok := key.(lua.LString); ok {
			if strValue, ok := value.(lua.LString); ok {
				result[string(strKey)] = string(strValue)
			}
		}
	})
	return result
}

// luaSearchItemFromTable converts a Lua table to LuaSearchItem
func luaSearchItemFromTable(L *lua.LState, t *lua.LTable) *LuaSearchItem {
	item := &LuaSearchItem{}

	if name := t.RawGetString("name"); name.Type() != lua.LTNil {
		item.Name = lua.LVAsString(name)
	}
	if key := t.RawGetString("key"); key.Type() != lua.LTNil {
		item.Key = lua.LVAsString(key)
	}
	if url := t.RawGetString("url"); url.Type() != lua.LTNil {
		item.URL = lua.LVAsString(url)
	}
	if summary := t.RawGetString("summary"); summary.Type() != lua.LTNil {
		item.Summary = lua.LVAsString(summary)
	}
	if cover := t.RawGetString("cover"); cover.Type() != lua.LTNil {
		item.Cover = lua.LVAsString(cover)
	}
	if scraperName := t.RawGetString("scraper_name"); scraperName.Type() != lua.LTNil {
		item.ScraperName = lua.LVAsString(scraperName)
	}

	return item
}

// luaGameItemFromTable converts a Lua table to LuaGameItem
func luaGameItemFromTable(L *lua.LState, t *lua.LTable) *LuaGameItem {
	item := &LuaGameItem{}

	if vndbID := t.RawGetString("vndb_id"); vndbID != lua.LNil {
		item.VNDBID = lua.LVAsString(vndbID)
	}
	if janCode := t.RawGetString("jan_code"); janCode != lua.LNil {
		item.JanCode = lua.LVAsString(janCode)
	}
	if dlCode := t.RawGetString("dl_code"); dlCode != lua.LNil {
		item.DlCode = lua.LVAsString(dlCode)
	}
	if name := t.RawGetString("name"); name != lua.LNil {
		item.Name = lua.LVAsString(name)
	}
	if cover := t.RawGetString("cover"); cover != lua.LNil {
		item.Cover = lua.LVAsString(cover)
	}
	if price := t.RawGetString("price"); price != lua.LNil {
		item.Price = lua.LVAsString(price)
	}
	if issueDate := t.RawGetString("issue_date"); issueDate != lua.LNil {
		item.IssueDate = lua.LVAsString(issueDate)
	}
	if story := t.RawGetString("story"); story != lua.LNil {
		item.Story = lua.LVAsString(story)
	}
	if platform := t.RawGetString("platform"); platform != lua.LNil {
		item.Platform = lua.LVAsString(platform)
	}
	if otherInfo := t.RawGetString("other_info"); otherInfo != lua.LNil {
		item.OtherInfo = lua.LVAsString(otherInfo)
	}

	// Parse arrays
	if aliasTbl, ok := t.RawGetString("alias").(*lua.LTable); ok {
		aliasTbl.ForEach(func(_, value lua.LValue) {
			if strValue, ok := value.(lua.LString); ok {
				item.Alias = append(item.Alias, string(strValue))
			}
		})
	}

	if imagesTbl, ok := t.RawGetString("images").(*lua.LTable); ok {
		imagesTbl.ForEach(func(_, value lua.LValue) {
			if strValue, ok := value.(lua.LString); ok {
				item.Images = append(item.Images, string(strValue))
			}
		})
	}

	// Parse category
	if categoryTbl, ok := t.RawGetString("category").(*lua.LTable); ok {
		item.Category = &LuaCategory{Name: lua.LVAsString(categoryTbl.RawGetString("name"))}
	}

	// Parse series
	if seriesTbl, ok := t.RawGetString("series").(*lua.LTable); ok {
		seriesTbl.ForEach(func(_, value lua.LValue) {
			if tbl, ok := value.(*lua.LTable); ok {
				item.Series = append(item.Series, &LuaSeries{Name: lua.LVAsString(tbl.RawGetString("name"))})
			}
		})
	}

	// Parse brands
	if brandsTbl, ok := t.RawGetString("brands").(*lua.LTable); ok {
		brandsTbl.ForEach(func(_, value lua.LValue) {
			if tbl, ok := value.(*lua.LTable); ok {
				brand := &LuaBrand{
					VNDBID: lua.LVAsString(tbl.RawGetString("vndb_id")),
					Name:   lua.LVAsString(tbl.RawGetString("name")),
				}
				if linksTbl, ok := tbl.RawGetString("links").(*lua.LTable); ok {
					linksTbl.ForEach(func(_, v lua.LValue) {
						if linkTbl, ok := v.(*lua.LTable); ok {
							brand.Links = append(brand.Links, &LuaLink{
								Name: lua.LVAsString(linkTbl.RawGetString("name")),
								Url:  lua.LVAsString(linkTbl.RawGetString("url")),
								Type: lua.LVAsString(linkTbl.RawGetString("type")),
							})
						}
					})
				}
				item.Brands = append(item.Brands, brand)
			}
		})
	}

	// Parse tags
	if tagsTbl, ok := t.RawGetString("tags").(*lua.LTable); ok {
		tagsTbl.ForEach(func(_, value lua.LValue) {
			if tbl, ok := value.(*lua.LTable); ok {
				item.Tags = append(item.Tags, &LuaTag{
					Name: lua.LVAsString(tbl.RawGetString("name")),
					Lang: lua.LVAsString(tbl.RawGetString("lang")),
				})
			}
		})
	}

	// Parse characters
	if charactersTbl, ok := t.RawGetString("characters").(*lua.LTable); ok {
		charactersTbl.ForEach(func(_, value lua.LValue) {
			if tbl, ok := value.(*lua.LTable); ok {
				char := &LuaCharacter{
					VNDBID:   lua.LVAsString(tbl.RawGetString("vndb_id")),
					Name:     lua.LVAsString(tbl.RawGetString("name")),
					Gender:   lua.LVAsString(tbl.RawGetString("gender")),
					Relation: lua.LVAsString(tbl.RawGetString("relation")),
					Summary:  lua.LVAsString(tbl.RawGetString("summary")),
					Cover:    lua.LVAsString(tbl.RawGetString("cover")),
				}

				if aliasTbl, ok := tbl.RawGetString("alias").(*lua.LTable); ok {
					aliasTbl.ForEach(func(_, v lua.LValue) {
						if strValue, ok := v.(lua.LString); ok {
							char.Alias = append(char.Alias, string(strValue))
						}
					})
				}

				if imagesTbl, ok := tbl.RawGetString("images").(*lua.LTable); ok {
					imagesTbl.ForEach(func(_, v lua.LValue) {
						if strValue, ok := v.(lua.LString); ok {
							char.Images = append(char.Images, string(strValue))
						}
					})
				}

				if cvTbl, ok := tbl.RawGetString("cv").(*lua.LTable); ok {
					char.CV = &LuaStaff{
						VNDBID:  lua.LVAsString(cvTbl.RawGetString("vndb_id")),
						Name:    lua.LVAsString(cvTbl.RawGetString("name")),
						Summary: lua.LVAsString(cvTbl.RawGetString("summary")),
						Gender:  lua.LVAsString(cvTbl.RawGetString("gender")),
						Cover:   lua.LVAsString(cvTbl.RawGetString("cover")),
					}
					if aliasTbl, ok := cvTbl.RawGetString("alias").(*lua.LTable); ok {
						aliasTbl.ForEach(func(_, v lua.LValue) {
							if strValue, ok := v.(lua.LString); ok {
								char.CV.Alias = append(char.CV.Alias, string(strValue))
							}
						})
					}
				}

				if tagsTbl, ok := tbl.RawGetString("tags").(*lua.LTable); ok {
					tagsTbl.ForEach(func(_, v lua.LValue) {
						if tagTbl, ok := v.(*lua.LTable); ok {
							char.Tags = append(char.Tags, &LuaTag{
								Name: lua.LVAsString(tagTbl.RawGetString("name")),
								Lang: lua.LVAsString(tagTbl.RawGetString("lang")),
							})
						}
					})
				}

				item.Characters = append(item.Characters, char)
			}
		})
	}

	// Parse staff
	if staffTbl, ok := t.RawGetString("staff").(*lua.LTable); ok {
		staffTbl.ForEach(func(_, value lua.LValue) {
			if tbl, ok := value.(*lua.LTable); ok {
				staff := &LuaStaff{
					VNDBID:  lua.LVAsString(tbl.RawGetString("vndb_id")),
					Name:    lua.LVAsString(tbl.RawGetString("name")),
					Summary: lua.LVAsString(tbl.RawGetString("summary")),
					Gender:  lua.LVAsString(tbl.RawGetString("gender")),
					Cover:   lua.LVAsString(tbl.RawGetString("cover")),
				}

				if aliasTbl, ok := tbl.RawGetString("alias").(*lua.LTable); ok {
					aliasTbl.ForEach(func(_, v lua.LValue) {
						if strValue, ok := v.(lua.LString); ok {
							staff.Alias = append(staff.Alias, string(strValue))
						}
					})
				}

				if imagesTbl, ok := tbl.RawGetString("images").(*lua.LTable); ok {
					imagesTbl.ForEach(func(_, v lua.LValue) {
						if strValue, ok := v.(lua.LString); ok {
							staff.Images = append(staff.Images, string(strValue))
						}
					})
				}

				if relationTbl, ok := tbl.RawGetString("relation").(*lua.LTable); ok {
					relationTbl.ForEach(func(_, v lua.LValue) {
						if strValue, ok := v.(lua.LString); ok {
							staff.Relation = append(staff.Relation, string(strValue))
						}
					})
				}

				if tagsTbl, ok := tbl.RawGetString("tags").(*lua.LTable); ok {
					tagsTbl.ForEach(func(_, v lua.LValue) {
						if tagTbl, ok := v.(*lua.LTable); ok {
							staff.Tags = append(staff.Tags, &LuaTag{
								Name: lua.LVAsString(tagTbl.RawGetString("name")),
								Lang: lua.LVAsString(tagTbl.RawGetString("lang")),
							})
						}
					})
				}

				item.Staff = append(item.Staff, staff)
			}
		})
	}

	// Parse links
	if linksTbl, ok := t.RawGetString("links").(*lua.LTable); ok {
		linksTbl.ForEach(func(_, value lua.LValue) {
			if tbl, ok := value.(*lua.LTable); ok {
				item.Links = append(item.Links, &LuaLink{
					Name: lua.LVAsString(tbl.RawGetString("name")),
					Url:  lua.LVAsString(tbl.RawGetString("url")),
					Type: lua.LVAsString(tbl.RawGetString("type")),
				})
			}
		})
	}

	return item
}

// ToSearchItem converts LuaSearchItem to scraper.SearchItem
func (l *LuaSearchItem) ToSearchItem(scraperName string) *scraper.SearchItem {
	return &scraper.SearchItem{
		Name:        l.Name,
		Key:         l.Key,
		URl:         l.URL,
		Summary:     l.Summary,
		Cover:       l.Cover,
		ScraperName: scraperName,
	}
}

// ToGameItem converts LuaGameItem to scraper.GameItem
func (l *LuaGameItem) ToGameItem(scraperName string) *scraper.GameItem {
	gameItem := &scraper.GameItem{
		GameVo: handler.GameVo{
			VNDBID:     l.VNDBID,
			JanCode:    l.JanCode,
			DlCode:     l.DlCode,
			Name:       l.Name,
			Alias:      l.Alias,
			Cover:      l.Cover,
			Images:     l.Images,
			Price:      l.Price,
			Story:      l.Story,
			Platform:   l.Platform,
			OtherInfo:  l.OtherInfo,
			Links:      make([]model.Link, 0),
		},
		ScraperName: scraperName,
	}

	if l.Category != nil {
		gameItem.GameVo.Category = &model.Category{Name: l.Category.Name}
	}

	for _, s := range l.Series {
		gameItem.GameVo.Series = append(gameItem.GameVo.Series, &model.Series{Name: s.Name})
	}

	for _, b := range l.Brands {
		brand := &model.Brand{
			VNDBID: b.VNDBID,
			Name:   b.Name,
			Links:  make([]model.Link, 0),
		}
		for _, link := range b.Links {
			brand.Links = append(brand.Links, model.Link{
				Name: link.Name,
				Url:  link.Url,
				Type: model.LinkType(link.Type),
			})
		}
		gameItem.GameVo.Brands = append(gameItem.GameVo.Brands, brand)
	}

	for _, t := range l.Tags {
		gameItem.GameVo.Tags = append(gameItem.GameVo.Tags, &model.Tag{
			Name: t.Name,
			Lang: t.Lang,
		})
	}

	for _, c := range l.Characters {
		char := handler.CharacterVo{
			VNDBID:   c.VNDBID,
			Name:     c.Name,
			Alias:    c.Alias,
			Gender:   model.Gender(c.Gender),
			Rlation:  model.CharacterRelation(c.Relation),
			Summary:  c.Summary,
			Cover:    c.Cover,
			Images:   c.Images,
			Tags:     make([]model.Tag, 0),
			PersonalInfo: model.PersonalInfo{},
		}
		for _, t := range c.Tags {
			char.Tags = append(char.Tags, model.Tag{
				Name: t.Name,
				Lang: t.Lang,
			})
		}
		if c.CV != nil {
			char.CV = handler.StaffVo{
				VNDBID:  c.CV.VNDBID,
				Name:    c.CV.Name,
				Alias:   c.CV.Alias,
				Summary: c.CV.Summary,
				Gender:  model.Gender(c.CV.Gender),
				Cover:   c.CV.Cover,
				Images:  c.CV.Images,
				Links:   make([]model.Link, 0),
			}
		}
		gameItem.GameVo.Characters = append(gameItem.GameVo.Characters, char)
	}

	for _, s := range l.Staff {
		staff := handler.StaffVo{
			VNDBID:  s.VNDBID,
			Name:    s.Name,
			Alias:   s.Alias,
			Summary: s.Summary,
			Gender:  model.Gender(s.Gender),
			Cover:   s.Cover,
			Images:  s.Images,
			Tags:    make([]model.Tag, 0),
			Links:   make([]model.Link, 0),
		}
		for _, r := range s.Relation {
			staff.Relation = append(staff.Relation, model.PersonRelation(r))
		}
		for _, t := range s.Tags {
			staff.Tags = append(staff.Tags, model.Tag{
				Name: t.Name,
				Lang: t.Lang,
			})
		}
		gameItem.GameVo.Staff = append(gameItem.GameVo.Staff, staff)
	}

	for _, link := range l.Links {
		gameItem.GameVo.Links = append(gameItem.GameVo.Links, model.Link{
			Name: link.Name,
			Url:  link.Url,
			Type: model.LinkType(link.Type),
		})
	}

	return gameItem
}