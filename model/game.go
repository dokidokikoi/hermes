package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type LinkType string

const (
	LinkTypeMv          LinkType = "mv"           // mv
	LinkTypeInfo        LinkType = "info"         // 介绍
	LinkTypeWorkThrough LinkType = "work_through" // 攻略
)

var LinkTypeMap = map[string]LinkType{
	"video": LinkTypeMv,
}

type Link struct {
	Name string   `json:"name"`
	Url  string   `json:"url"`
	Type LinkType `json:"type"`
}

type Links []Link

func (a *Links) scanBytes(src []byte) error {
	return json.Unmarshal(src, a)
}

// Scan implements the sql.Scanner interface.
func (a *Links) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("cannot convert %T to Link", src)
}

// Value implements the driver.Valuer interface.
func (a Links) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	data, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type DownloadType string

const (
	DownloadMagnet = "magnet"
)

type DownloadInfo struct {
	Content string       `json:"content"`
	Type    DownloadType `json:"type"`
	Version string       `json:"version"`
}

type DownloadInfos []DownloadInfo

// Scan implements the sql.Scanner interface.
func (a *DownloadInfos) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, a)
	case string:
		return json.Unmarshal([]byte(src), a)
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("cannot convert %T to Link", src)
}

// Value implements the driver.Valuer interface.
func (a DownloadInfos) Value() (driver.Value, error) {
	if a == nil {
		return "[]", nil
	}

	data, err := json.Marshal(a)
	if err != nil {
		return "[]", err
	}

	return data, nil
}

type GameInstance struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	GameID    uint      `gorm:"uniqueIndex:uk_game_version"`
	Version   string    `gorm:"uniqueIndex:uk_game_version" json:"version"`
	Path      string    `gorm:"index:idx_game_path" json:"path"`
	Size      int64     `json:"size"`
	Language  string    `json:"language"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
}

func (GameInstance) TableName() string {
	return "game_instances"
}

type Game struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	JanCode     string         `json:"jan_code"`
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	Alias       pq.StringArray `gorm:"type:text[]" json:"alias"`
	Cover       string         `json:"cover"`
	Images      pq.StringArray `gorm:"type:text[]" json:"images"`
	Category    *Category      `gorm:"foreignKey:CategoryID" json:"category"`
	CategoryID  uint           `gorm:"default:null" json:"-"`
	Series      []Series       `gorm:"many2many:game_series;" json:"series"`
	Developer   *Developer     `gorm:"foreignKey:DeveloperID" json:"developer"`
	DeveloperID uint           `gorm:"default:null" json:"-"`
	Publisher   *Publisher     `gorm:"foreignKey:PublisherID" json:"publisher"`
	PublisherID uint           `gorm:"default:null" json:"-"`
	Price       string         `json:"price"`
	IssueDate   time.Time      `json:"issue_date"`
	Story       string         `json:"story"`
	Platform    string         `json:"platform"`
	Tags        []Tag          `gorm:"many2many:game_tag;" json:"tags"`
	Links       Links          `gorm:"type:json" json:"links"`
	OtherInfo   string         `json:"other_info"`
	CreatedAt   time.Time      `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime:milli"`
}

func (Game) TableName() string {
	return "games"
}

type GameSeries struct {
	GameID   uint `gorm:"primaryKey"`
	SeriesID uint `gorm:"primaryKey"`
}

func (GameSeries) TableName() string {
	return "game_series"
}

type CharacterRelation string

const (
	CRelationMain  CharacterRelation = "main"
	CRelationMinor CharacterRelation = "minor"
	CRelationMob   CharacterRelation = "mob"
)

type GameCharacter struct {
	GameID      uint              `gorm:"primaryKey"`
	CharacterID uint              `gorm:"primaryKey"`
	Relation    CharacterRelation `json:"relation"`
	Character   *Character        `gorm:"-" json:"-"`
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

type PersonRelation string

const (
	PRelationUnknown PersonRelation = "unknown"
	PRelationPainter PersonRelation = "painter"
	PRelationWriter  PersonRelation = "writer"
	PRelationCV      PersonRelation = "cv"
	PRelationMusic   PersonRelation = "music"
)

type GameStaff struct {
	GameID    uint             `gorm:"primaryKey"`
	PersonID  uint             `gorm:"primaryKey"`
	Relation  PersonRelation   `json:"relation"`
	Relations []PersonRelation `gorm:"-" json:"-"`
	Person    *Person          `gorm:"-" json:"-"`
}

func (GameStaff) TableName() string {
	return "game_staff"
}
