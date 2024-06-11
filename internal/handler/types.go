package handler

import (
	"hermes/model"
	"time"
)

type GameListReq struct {
	Keyword        string      `json:"keyword"`
	Tags           []uint      `json:"tags"`
	Series         uint        `json:"series"`
	Category       uint        `json:"category"`
	Staff          uint        `json:"staff"`
	Character      uint        `json:"character"`
	Developer      uint        `json:"developer"`
	Publisher      uint        `json:"publisher"`
	SizeRange      []int64     `json:"size_range"`
	IssueDateRange []time.Time `json:"issue_date_range"`
	CreatedAtRange []time.Time `json:"created_at_range"`
	FullText       bool        `json:"full_text"`
	Page           int         `json:"page"`
	PageSize       int         `json:"page_size"`
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
	ID         uint             `json:"id"`
	JanCode    string           `json:"jan_code"`
	Code       string           `json:"code"`
	Name       string           `json:"name"`
	Alias      []string         `json:"alias"`
	Cover      string           `json:"cover"`
	Images     []string         `json:"images"`
	Versions   []string         `json:"versions"`
	Category   *model.Category  `json:"category"`
	Series     []model.Series   `json:"series"`
	Developer  *model.Developer `json:"developer"`
	Publisher  *model.Publisher `json:"publisher"`
	Price      string           `json:"price"`
	IssueDate  time.Time        `json:"issue_date"`
	Story      string           `json:"story"`
	Platform   string           `json:"platform"`
	Tags       []model.Tag      `json:"tags"`
	Characters []CharacterVo    `json:"characters"`
	Links      []model.Link     `json:"links"`
	OtherInfo  string           `json:"other_info"`
	Staff      []StaffVo        `json:"staff"`
	CreatedAt  time.Time        `json:"created_at"`
}

type GameInsVo struct {
	ID         uint             `json:"id"`
	Name       string           `json:"name"`
	Alias      []string         `json:"alias"`
	Cover      string           `json:"cover"`
	Images     []string         `json:"images"`
	Version    string           `json:"version"`
	Category   *model.Category  `json:"category"`
	Series     []model.Series   `json:"series"`
	Developer  *model.Developer `json:"developer"`
	Publisher  *model.Publisher `json:"publisher"`
	Price      string           `json:"price"`
	IssueDate  time.Time        `json:"issue_date"`
	Story      string           `json:"story"`
	Platform   string           `json:"platform"`
	Tags       []model.Tag      `json:"tags"`
	Characters []CharacterVo    `json:"characters"`
	Links      []model.Link     `json:"links"`
	OtherInfo  string           `json:"other_info"`
	Staff      []StaffVo        `json:"staff"`
	CreatedAt  time.Time        `json:"created_at"`

	Size     int64  `json:"size"`
	Path     string `json:"path"`
	Language string `json:"language"`
	Comment  string `json:"comment"`
}

type StaffVo struct {
	ID        uint        `json:"id"`
	Name      string      `json:"name"`
	Alias     []string    `json:"alias"`
	Cover     string      `json:"cover"`
	Images    []string    `json:"images"`
	Tags      []model.Tag `json:"tags"`
	Summary   string      `json:"summary"`
	Gender    string      `json:"gender"`
	Relation  []string    `json:"relation"`
	CreatedAt time.Time   `json:"created_at"`
}

type CharacterVo struct {
	ID        uint        `json:"id"`
	Name      string      `json:"name"`
	Alias     []string    `json:"alias"`
	Gender    string      `json:"gender"`
	Rlation   string      `json:"relation"`
	Summary   string      `json:"summary"`
	Cover     string      `json:"cover"`
	Images    []string    `json:"images"`
	CV        StaffVo     `json:"cv"`
	Tags      []model.Tag `json:"tags"`
	CreatedAt time.Time   `json:"created_at"`
}

type CharacterListReq struct {
	Keyword        string       `json:"keyword"`
	Tags           []uint       `json:"tags"`
	CreatedAtRange []time.Time  `json:"created_at_range"`
	Gender         model.Gender `json:"gender"`
	CV             uint         `json:"cv"`
	FullText       bool         `json:"full_text"`
}

type PersonListReq struct {
	Keyword        string       `json:"keyword"`
	Tags           []uint       `json:"tags"`
	CreatedAtRange []time.Time  `json:"created_at_range"`
	Gender         model.Gender `json:"gender"`
	CV             uint         `json:"cv"`
	FullText       bool         `json:"full_text"`
}

type ScraperSearchReq struct {
	Name      string `json:"name"`
	Keyword   string `json:"keyword"`
	Page      int    `json:"page"`
	RequestID string `json:"request_id"`
}

type ScraperDetailReq struct {
	RequestID string `json:"request_id"`
	Path      string `json:"path"`
	Version   string `json:"version"`
	Objs      []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"objs"`
}

type ScraperGetReq struct {
	RequestID string `form:"request_id"`
}
