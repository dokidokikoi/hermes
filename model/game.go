package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type LinkType string

const (
	LinkTypeMv          LinkType = "mv"           // mv
	LinkTypeInfo        LinkType = "info"         // 介绍
	LinkTypeWorkThrough LinkType = "work_through" // 攻略
	LinkTypeOfficial    LinkType = "official"
	LinkTypeOther       LinkType = "other"
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
	Version   string    `gorm:"uniqueIndex:uk_game_version;type:verchar(100)" json:"version"`
	Path      string    `gorm:"index:idx_game_path;type:verchar(255)" json:"path"`
	Size      int64     `json:"size"`
	Language  Array     `gorm:"type:json" json:"language"`
	Comment   string    `json:"comment"`
	Platform  Array     `gorm:"type:json" json:"platform"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
}

func (GameInstance) TableName() string {
	return "game_instances"
}

type Game struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UUID       string    `gorm:"uniqueIndex;type:varchar(100)" json:"uuid"`
	VNDBID     string    `gorm:"type:varchar(255);index:game_vndbid_idx" json:"vndb_id"`
	JanCode    string    `gorm:"type:varchar(32)" json:"jan_code"`
	DlCode     string    `gorm:"type:varchar(32)" json:"dl_code"`
	Name       string    `gorm:"type:varchar(255);index:game_name_idx" json:"name"`
	Alias      Array     `gorm:"type:json" json:"alias"`
	Cover      string    `gorm:"type:varchar(512)" json:"cover"`
	Images     Array     `gorm:"type:json" json:"images"`
	Category   *Category `gorm:"foreignKey:CategoryID" json:"category"`
	CategoryID uint      `gorm:"default:null" json:"-"`
	Series     []*Series `gorm:"many2many:game_series;" json:"series"`
	Brands     []*Brand  `gorm:"many2many:game_brands;" json:"brands"`
	Price      string    `gorm:"type:varchar(32)" json:"price"`
	IssueDate  time.Time `json:"issue_date"`
	Story      string    `json:"story"`
	Tags       []*Tag    `gorm:"many2many:game_tag;" json:"tags"`
	Links      Links     `gorm:"type:json" json:"links"`
	OtherInfo  string    `json:"other_info"`
	CreatedAt  time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime:milli"`
}

func (Game) TableName() string {
	return "games"
}

type GameBrands struct {
	GameID  uint `gorm:"primaryKey"`
	BrandID uint `gorm:"primaryKey"`
}

func (GameBrands) TableName() string {
	return "game_brands"
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
	Relation    CharacterRelation `gorm:"type:varchar(32)" json:"relation"`
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

type PersonRelations []PersonRelation

func (a *PersonRelations) scanBytes(src []byte) error {
	return json.Unmarshal(src, a)
}

// Scan implements the sql.Scanner interface.
func (a *PersonRelations) Scan(src interface{}) error {
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
func (a PersonRelations) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	data, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	return data, nil
}

const (
	PRelationUnknown PersonRelation = "unknown"
	PRelationPainter PersonRelation = "painter"
	PRelationWriter  PersonRelation = "writer"
	PRelationCV      PersonRelation = "cv"
	PRelationMusic   PersonRelation = "music"
)

type GameStaff struct {
	GameID    uint            `gorm:"primaryKey"`
	PersonID  uint            `gorm:"primaryKey"`
	Relations PersonRelations `gorm:"type:json"`
	Person    *Person         `gorm:"-" json:"-"`
}

func (GameStaff) TableName() string {
	return "game_staff"
}
