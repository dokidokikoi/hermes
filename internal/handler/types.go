package handler

import (
	"izumi/model"
	"strings"
	"time"
)

type GameListReq struct {
	Keyword        string      `json:"keyword"`
	TagIDs         []uint      `json:"tags"`
	Series         uint        `json:"series"`
	Category       uint        `json:"category"`
	Staff          uint        `json:"staff"`
	Character      uint        `json:"character"`
	Brand          uint        `json:"brand"`
	SizeRange      []int64     `json:"size_range"`
	IssueDateRange []time.Time `json:"issue_date_range"`
	CreatedAtRange []time.Time `json:"created_at_range"`
	FullText       bool        `json:"full_text"`
	Page           int         `json:"page"`
	PageSize       int         `json:"page_size"`
	OrderBy        string      `json:"order_by"`
}

type GameInlineVo struct {
	ID        uint      `json:"id"`
	JanCode   string    `json:"jan_code"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Alias     []string  `json:"alias"`
	Cover     string    `json:"cover"`
	IssueDate time.Time `json:"issue_date"`
	Story     string    `json:"story"`
}

type GameVo struct {
	ID            uint                `json:"id"`
	RelIDs        []string            `json:"rel_ids"`
	JanCode       string              `json:"jan_code"`
	DlCode        string              `json:"dl_code"`
	Name          string              `json:"name"`
	Alias         []string            `json:"alias"`
	Cover         string              `json:"cover"`
	Images        []string            `json:"images"`
	Category      *model.Category     `json:"category"`
	Series        []*model.Series     `json:"series"`
	Brands        []*model.Brand      `json:"brands"`
	Price         string              `json:"price"`
	IssueDate     time.Time           `json:"issue_date"`
	Story         string              `json:"story"`
	Platform      string              `json:"platform"`
	Tags          []*model.Tag        `json:"tags"`
	Characters    []CharacterVo       `json:"characters"`
	Links         []model.Link        `json:"links"`
	OtherInfo     string              `json:"other_info"`
	Staff         []StaffVo           `json:"staff"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	DownloadInfos model.DownloadInfos `json:"download_info"`

	// 仅用于 info 文件
	Version  string      `json:"version,omitempty"`
	Size     int64       `json:"size,omitempty"`
	Language model.Array `json:"language,omitempty"`
	Comment  string      `json:"comment,omitempty"`
}

func Vo2Game(game GameVo) model.Game {
	return model.Game{
		ID:       game.ID,
		RelIDs:   game.RelIDs,
		Name:     strings.TrimSpace(game.Name),
		Alias:    game.Alias,
		Cover:    game.Cover,
		Images:   game.Images,
		Category: game.Category,
		CategoryID: func() uint {
			if game.Category != nil {
				return game.Category.ID
			}
			return 0
		}(),
		Series:    game.Series,
		Brands:    game.Brands,
		Price:     game.Price,
		IssueDate: game.IssueDate,
		Story:     strings.TrimSpace(game.Story),
		Tags:      game.Tags,
		Links:     game.Links,
		OtherInfo: strings.TrimSpace(game.OtherInfo),
	}
}

func Game2Vo(game model.Game, characters []CharacterVo, staff []StaffVo) GameVo {
	return GameVo{
		ID:         game.ID,
		Name:       game.Name,
		Alias:      game.Alias,
		Cover:      game.Cover,
		Images:     game.Images,
		Category:   game.Category,
		Series:     game.Series,
		Brands:     game.Brands,
		Price:      game.Price,
		IssueDate:  game.IssueDate,
		Story:      game.Story,
		Tags:       game.Tags,
		Characters: characters,
		Links:      game.Links,
		Staff:      staff,
		CreatedAt:  game.CreatedAt,
		UpdatedAt:  game.UpdatedAt,
	}
}

type StaffVo struct {
	ID        uint                   `json:"id"`
	RelIDs    []string               `json:"rel_ids"`
	Name      string                 `json:"name"`
	Alias     []string               `json:"alias"`
	Cover     string                 `json:"cover"`
	Images    []string               `json:"images"`
	Tags      []model.Tag            `json:"tags"`
	Summary   string                 `json:"summary"`
	Gender    model.Gender           `json:"gender"`
	Relation  []model.PersonRelation `json:"relation"`
	Links     []model.Link           `json:"links"`
	CreatedAt time.Time              `json:"created_at"`
}

func Vo2Person(staff StaffVo) model.Person {
	return model.Person{
		ID:      staff.ID,
		Name:    staff.Name,
		Alias:   staff.Alias,
		Cover:   staff.Cover,
		Images:  staff.Images,
		Summary: staff.Summary,
		Gender:  staff.Gender,
		Links:   staff.Links,
	}
}

func Person2Vo(person model.Person, relation []model.PersonRelation) StaffVo {
	return StaffVo{
		ID:       person.ID,
		Name:     person.Name,
		Alias:    person.Alias,
		Cover:    person.Cover,
		Images:   person.Images,
		Summary:  person.Summary,
		Gender:   person.Gender,
		Links:    person.Links,
		Relation: relation,
	}
}

type CharacterVo struct {
	ID           uint                    `json:"id"`
	RelIDs       []string                `json:"rel_ids"`
	Name         string                  `json:"name"`
	Alias        []string                `json:"alias"`
	Gender       model.Gender            `json:"gender"`
	Rlation      model.CharacterRelation `json:"relation"`
	Summary      string                  `json:"summary"`
	Cover        string                  `json:"cover"`
	Images       []string                `json:"images"`
	CV           StaffVo                 `json:"cv"`
	Tags         []model.Tag             `json:"tags"`
	CreatedAt    time.Time               `json:"created_at"`
	Games        []CharacterGameVo       `json:"games,omitempty"`
	Weight       int8                    `json:"weight"`
	PersonalInfo model.PersonalInfo      `json:"personal_info"`
}

func Vo2Character(character CharacterVo) model.Character {
	return model.Character{
		ID:           character.ID,
		Name:         character.Name,
		Alias:        character.Alias,
		Gender:       character.Gender,
		Summary:      character.Summary,
		Cover:        character.Cover,
		Images:       character.Images,
		CV:           Vo2Person(character.CV),
		PersonID:     character.CV.ID,
		Weight:       character.Weight,
		PersonalInfo: character.PersonalInfo,
	}
}

func Character2Vo(character model.Character, relation model.CharacterRelation, games []CharacterGameVo) CharacterVo {
	return CharacterVo{
		ID:           character.ID,
		Name:         character.Name,
		Alias:        character.Alias,
		Gender:       character.Gender,
		Rlation:      relation,
		Summary:      character.Summary,
		Cover:        character.Cover,
		Images:       character.Images,
		CV:           Person2Vo(character.CV, nil),
		CreatedAt:    character.CreatedAt,
		Weight:       character.Weight,
		PersonalInfo: character.PersonalInfo,
		Games:        games,
	}
}

type CharacterGameVo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type PersonListReq struct {
	Keyword        string       `json:"keyword"`
	Tags           []uint       `json:"tags"`
	CreatedAtRange []time.Time  `json:"created_at_range"`
	Gender         model.Gender `json:"gender"`
	FullText       bool         `json:"full_text"`

	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	OrderBy  string `json:"order_by"`
}

type CharacterSearchReq struct {
	Keyword        string       `json:"keyword"`
	Tags           []uint       `json:"tags"`
	CreatedAtRange []time.Time  `json:"created_at_range"`
	Gender         model.Gender `json:"gender"`
	CV             uint         `json:"cv"`
	FullText       bool         `json:"full_text"`

	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	OrderBy  string `json:"order_by"`
}

type ScraperSearchReq struct {
	Name      string `json:"name"`
	Keyword   string `json:"keyword"`
	Page      int    `json:"page"`
	RequestID string `json:"request_id"`
}

type ScraperDetailReq struct {
	RequestID string           `json:"request_id"`
	Objs      []model.ScrapObj `json:"objs"`
}

type ScraperGetReq struct {
	RequestID string `form:"request_id"`
}

type UpdateProxyReq struct {
	Key    string `json:"key" binding:"required"`
	Policy string `json:"policy" binding:"required"`
}

type LsReq struct {
	Path        string `form:"path"`
	OnlyNoScrap bool   `form:"only_no_scrap"`
}

type SeriesVo struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Games     int64     `json:"games"`
	CreatedAt time.Time `json:"created_at"`
}
