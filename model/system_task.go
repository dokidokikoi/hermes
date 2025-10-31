package model

import "time"

type SystemTaskType string

const (
	SystemTaskTypeDownload SystemTaskType = "download"
	SystemTaskTypeLoad     SystemTaskType = "load"
)

type SystemTaskState string

const (
	SystemTaskStateRunning  SystemTaskState = "running"
	SystemTaskStateDone     SystemTaskState = "done"
	SystemTaskStateFailed   SystemTaskState = "failed"
	SystemTaskStateCanceled SystemTaskState = "canceled"
)

type SystemTask struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Amount    int             `json:"amount"`
	Type      SystemTaskType  `gorm:"type:varchar" json:"type"`
	State     SystemTaskState `gorm:"type:varchar" json:"state"`
	CreatedAt time.Time       `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime:milli"`
}

func (SystemTask) TableName() string {
	return "system_task"
}
