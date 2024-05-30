package model

import (
	"time"

	"github.com/lib/pq"
)

type Vedio struct {
	Title       string
	No          string
	Cover       string
	Images      pq.StringArray `gorm:"type:text[]"`
	Tags        []Tag          `gorm:"many2many:vedio_tag;"`
	Series      []Series       `gorm:"many2many:vedio_series;"`
	Publisher   Publisher      `gorm:"foreignKey:PublisherID"`
	PublisherID uint           `gorm:"default:null"`
	Actors      []Person       `gorm:"many2many:vedio_actor;"`
	Story       string
	Path        string
	Size        string
	OtherInfo   string
	IssueDate   time.Time
	CreatedAt   time.Time `gorm:"autoCreateTime:milli"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:milli"`
}

type VedioTag struct {
	VedioID uint `gorm:"primaryKey"`
	TagID   uint `gorm:"primaryKey"`
}

func (VedioTag) TableName() string {
	return "vedio_tag"
}

type VedioActor struct {
	VedioID uint `gorm:"primaryKey"`
	ActorID uint `gorm:"primaryKey"`
}

func (VedioActor) TableName() string {
	return "vedio_actor"
}

type VedioSeries struct {
	VedioID  uint `gorm:"primaryKey"`
	SeriesID uint `gorm:"primaryKey"`
}

func (VedioSeries) TableName() string {
	return "vedio_series"
}
