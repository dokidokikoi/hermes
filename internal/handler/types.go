package handler

import (
	"izumi/model"
	"time"
)

type GameListReq struct {
	Keyword        string      `json:"keyword"`
	Tags           []string    `json:"tags"`
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
	UUID          string              `json:"uuid"`
	JanCode       string              `json:"jan_code"`
	Code          string              `json:"code"`
	Name          string              `json:"name"`
	Alias         []string            `json:"alias"`
	Cover         string              `json:"cover"`
	Images        []string            `json:"images"`
	Category      *model.Category     `json:"category"`
	Series        []*model.Series     `json:"series"`
	Brand         *model.Brand        `json:"brand"`
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
	DownloadInfos model.DownloadInfos `json:"download_info"`

	// 仅用于 info 文件
	Version  string      `json:"version,omitempty"`
	Size     int64       `json:"size,omitempty"`
	Language model.Array `json:"language,omitempty"`
	Comment  string      `json:"comment,omitempty"`
}

type StaffVo struct {
	ID        uint                   `json:"id"`
	UUID      string                 `json:"uuid"`
	Name      string                 `json:"name"`
	Alias     []string               `json:"alias"`
	Cover     string                 `json:"cover"`
	Images    []string               `json:"images"`
	Tags      []model.Tag            `json:"tags"`
	Summary   string                 `json:"summary"`
	Gender    model.Gender           `json:"gender"`
	Relation  []model.PersonRelation `json:"relation"`
	CreatedAt time.Time              `json:"created_at"`
}

type CharacterVo struct {
	ID        uint                    `json:"id"`
	UUID      string                  `json:"uuid"`
	Name      string                  `json:"name"`
	Alias     []string                `json:"alias"`
	Gender    model.Gender            `json:"gender"`
	Rlation   model.CharacterRelation `json:"relation"`
	Summary   string                  `json:"summary"`
	Cover     string                  `json:"cover"`
	Images    []string                `json:"images"`
	CV        StaffVo                 `json:"cv"`
	Tags      []model.Tag             `json:"tags"`
	CreatedAt time.Time               `json:"created_at"`
	Games     []CharacterGameVo       `json:"games,omitempty"`
	Weight    int8                    `json:"weight"`
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
	RequestID string `json:"request_id"`
	Objs      []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"objs"`
}

type ScraperGetReq struct {
	RequestID string `form:"request_id"`
}

type UpdateProxyReq struct {
	Key    string `json:"key" binding:"required"`
	Policy string `json:"policy" binding:"required"`
}

type LsReq struct {
	Path string `form:"path"`
}
