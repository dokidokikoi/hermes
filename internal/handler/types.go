package handler

import (
	"hermes/model"
	"time"
)

type GameListReq struct {
	Keyword        string      `json:"keyword"`
	Tags           []uint      `json:"tags"`
	Series         []uint      `json:"series"`
	Category       uint        `json:"category"`
	Character      uint        `json:"character"`
	Developer      uint        `json:"developer"`
	Publisher      uint        `json:"publisher"`
	SizeRange      []int64     `json:"size_range"`
	IssueDateRange []time.Time `json:"issue_date_range"`
	CreatedAtRange []time.Time `json:"created_at_range"`
	FullText       bool        `json:"full_text"`
}

type GameVo struct {
	ID         uint             `json:"id"`
	Name       string           `json:"name"`
	Alias      []string         `json:"alias"`
	Cover      string           `json:"cover"`
	Images     []string         `json:"images"`
	Version    string           `json:"version"`
	Category   *model.Category  `json:"category"`
	Series     []model.Series   `json:"series"`
	Developer  *model.Developer `gjson:"developer"`
	Publisher  *model.Publisher `json:"publisher"`
	Price      string           `json:"price"`
	IssueDate  time.Time        `json:"issue_date"`
	Story      string           `json:"story"`
	Platform   string           `json:"platform"`
	Tags       []model.Tag      `json:"tags"`
	Characters []CharacterVo    `json:"characters"`
	Language   string           `json:"language"`
	Links      []model.Link     `json:"links"`
	Path       string           `json:"path"`
	Size       int64            `json:"size"`
	OtherInfo  string           `json:"other_info"`
	Staff      []StaffVo        `json:"staff"`
	CreatedAt  time.Time        `json:"created_at"`
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
