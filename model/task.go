package model

import "time"

const (
	TaskStatusWait int8 = iota
	TaskStatusSucceed
	TaskStatusFailed
	TaskStatusCanceled
)

const (
	TaskTypeSearch int8 = iota
	TaskTypeDetail
)

type Task struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	RequestID   string `gorm:"varchar(64);index:idx_request_id" json:"request_id"`
	ScraperName string `gorm:"type:varchar(32)" json:"scraper_name"`
	Type        int8
	Param       string    `json:"param"`
	Result      string    `json:"result"`
	Status      int8      `json:"status"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	Duration    int64     `json:"duration"`
	CreatedAt   time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:milli"`
	DeletedAt   *time.Time
}

func (Task) TableName() string {
	return "tasks"
}
