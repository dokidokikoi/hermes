package model

import "time"

type Link struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Game struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Name        string     `json:"name"`
	Alias       []string   `gorm:"type:text[]" json:"alias"`
	Cover       string     `json:"cover"`
	Images      []string   `gorm:"type:text[]" json:"images"`
	Version     string     `json:"version"`
	Category    *Category  `gorm:"foreignKey:CategoryID" json:"category"`
	CategoryID  uint       `gorm:"default:null" json:"-"`
	Series      []Series   `gorm:"many2many:game_series;" json:"series"`
	Developer   *Developer `gorm:"foreignKey:DeveloperID" json:"developer"`
	DeveloperID uint       `gorm:"default:null" json:"-"`
	Publisher   *Publisher `gorm:"foreignKey:PublisherID" json:"publisher"`
	PublisherID uint       `gorm:"default:null" json:"-"`
	Price       string     `json:"price"`
	IssueDate   time.Time  `json:"issue_date"`
	Story       string     `json:"story"`
	Platform    string     `json:"platform"`
	Tags        []Tag      `gorm:"many2many:game_tag;" json:"tags"`
	Language    string     `json:"language"`
	Links       []Link     `gorm:"type:json" json:"links"`
	Path        string     `json:"path"`
	Size        int64      `json:"size"`
	OtherInfo   string     `json:"other_info"`
	CreatedAt   time.Time  `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime:milli"`
}

type GameSeries struct {
	GameID   uint `gorm:"primaryKey"`
	SeriesID uint `gorm:"primaryKey"`
}

func (GameSeries) TableName() string {
	return "game_series"
}

type CharacterRelation int8

const (
	CRelationNone CharacterRelation = iota
	CRelationMain
	CRelationMinor
	CRelationMob
)

var CRelationMap = map[string]CharacterRelation{
	"主角": CRelationMain,
	"配角": CRelationMinor,
	"路人": CRelationMob,
}

func (r CharacterRelation) String() string {
	switch r {
	case 1:
		return "主角"
	case 2:
		return "配角"
	case 3:
		return "路人"
	default:
		return "none"
	}
}

func Str2CRelation(str string) CharacterRelation {
	switch str {
	case "主角":
		return CRelationMain
	case "配角":
		return CRelationMinor
	case "路人":
		return CRelationMob
	default:
		return CRelationNone
	}
}

type GameCharacter struct {
	GameID      uint   `gorm:"primaryKey"`
	CharacterID uint   `gorm:"primaryKey"`
	Relation    string `json:"relation"`
}

func (GameCharacter) TableName() string {
	return "game_character"
}

type GameTag struct {
	GameID uint `gorm:"primaryKey"`
	TagID  uint `gorm:"primaryKey"`
}

func (GameTag) TableName() string {
	return "game_tag"
}

type PersonRelation int8

const (
	PRelationNone PersonRelation = iota
	PRelationPainter
	PRelationWriter
	PRelationCV
	PRelationMusic
)

var PRelationMap = map[string]PersonRelation{
	"原画": PRelationPainter,
	"剧本": PRelationWriter,
	"声优": PRelationCV,
	"音乐": PRelationMusic,
}

func (r PersonRelation) String() string {
	switch r {
	case 1:
		return "原画"
	case 2:
		return "剧本"
	case 3:
		return "声优"
	case 4:
		return "音乐"
	default:
		return "none"
	}
}

type GameStaff struct {
	GameID   uint   `gorm:"primaryKey"`
	PersonID uint   `gorm:"primaryKey"`
	Relation string `json:"relation"`
}

func (GameStaff) TableName() string {
	return "game_staff"
}
