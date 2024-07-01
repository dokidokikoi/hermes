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
	RequestID   string `gorm:"index:idx_request_id" json:"request_id"`
	ScraperName string `json:"scraper_name"`
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

type RefGameInstance struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RequestID string    `gorm:"uniqueIndex:uni_request_id" json:"request_id"`
	Path      string    `json:"path"`
	Version   string    `json:"version"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
}

func (RefGameInstance) TableName() string {
	return "ref_gameInstances"
}
