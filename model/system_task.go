package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type SystemTaskType string

const (
	SystemTaskTypeDownload SystemTaskType = "download"
	SystemTaskTypeLoad     SystemTaskType = "load"
	SystemTaskTypeScrap    SystemTaskType = "scrap"
)

type SystemTaskState string

const (
	SystemTaskStateRunning  SystemTaskState = "running"
	SystemTaskStateDone     SystemTaskState = "done"
	SystemTaskStateFailed   SystemTaskState = "failed"
	SystemTaskStateCanceled SystemTaskState = "canceled"
)

type ScrapObj struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type SystemTaskParam struct {
	Amount    int        `json:"amount,omitempty"`
	ScrapObjs []ScrapObj `json:"objs"`
	Path      string     `json:"path"`
	Version   string     `json:"version"`
}

func (p *SystemTaskParam) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, p)
	case string:
		return json.Unmarshal([]byte(src), p)
	case nil:
		*p = SystemTaskParam{}
		return nil
	}

	return fmt.Errorf("cannot convert %T to Link", src)
}

func (p SystemTaskParam) Value() (driver.Value, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return "{}", err
	}

	return data, nil
}

type SystemTask struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Type      SystemTaskType  `gorm:"type:varchar" json:"type"`
	State     SystemTaskState `gorm:"type:varchar" json:"state"`
	Param     SystemTaskParam `gorm:"json" json:"param"`
	Result    string          `json:"result"`
	CreatedAt time.Time       `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime:milli"`
}

func (SystemTask) TableName() string {
	return "system_task"
}
