package model

import (
	"time"
)

const (
	NS_RECLASS   = "reclass"
	NS_LANUGAGE  = "language"
	NS_PARODY    = "parody"
	NS_CHARACTER = "character"
	NS_GROUP     = "group"
	NS_ARTIST    = "artist"
	NS_COSPLAYER = "cosplayer"
	NS_MALE      = "male"
	NS_FEMALE    = "female"
	NS_MIXED     = "mixed"
	NS_OTHER     = "other"
	NS_LOCATION  = "location"

	NS_RECLASS_ABBR   = "r"
	NS_LANUGAGE_ABBR  = "l"
	NS_PARODY_ABBR    = "p"
	NS_CHARACTER_ABBR = "c"
	NS_GROUP_ABBR     = "g"
	NS_ARTIST_ABBR    = "a"
	NS_COSPLAYER_ABBR = "cos"
	NS_MALE_ABBR      = "m"
	NS_FEMALE_ABBR    = "f"
	NS_MIXED_ABBR     = "x"
	NS_OTHER_ABBR     = "o"
	NS_LOCATION_ABBR  = "loc"
)

type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	NS        string    `grom:"type:varchar(16);uniqueIndex:uk_tag" json:"ns"`
	Key       string    `gorm:"type:varchar(64);uniqueIndex:uk_tag" json:"key"`
	Name      string    `gorm:"type:varchar(255);" json:"name"`
	Intro     string    `grom:"type:varchar(512);" json:"intro"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
}

func (Tag) TableName() string {
	return "tags"
}

type DecidedTag struct {
	ID    int64  `gorm:"type:varchar(32);primaryKey" json:"plat"`
	Tag   string `grom:"type:varchar(128);promaryKey" json:"rel_id"`
	TagID uint   `json:"tag_id"`
}

func (DecidedTag) TableName() string {
	return "decided_tags"
}
