package model

import (
	"time"

	"gorm.io/datatypes"
)

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
	ID          uint      `gorm:"primaryKey;" json:"id"`
	RequestID   string    `gorm:"type:varchar(64);index:idx_request_id" json:"request_id"`
	ScraperName string    `gorm:"type:varchar(32)" json:"scraper_name"`
	Type        int8      `gorm:"type:tinyint"`
	Param       string    `json:"param"`
	Result      string    `json:"result"`
	Status      int8      `json:"status"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	Duration    int64     `json:"duration"`
	CreatedAt   time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:milli"`
}

func (Task) TableName() string {
	return "tasks"
}

type UndecidedGameStatus string

const (
	UndecidedGameStatusWait     UndecidedGameStatus = "wait"
	UndecidedGameStatusAccepted UndecidedGameStatus = "accepted"
	UndecidedGameStatusRejected UndecidedGameStatus = "rejected"
)

type UndecidedGame struct {
	UUID    string                    `gorm:"type:varchar(100);primaryKey" json:"uuid"`
	Path    string                    `gorm:"type:varchar(255)" json:"path"`
	TaskIDs datatypes.JSONSlice[uint] `gorm:"type:jsonb" json:"task_ids"`
	Status  UndecidedGameStatus       `gorm:"type:varchar(32)" json:"status"`

	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
}

func (UndecidedGame) TableName() string {
	return "undecided_games"
}
