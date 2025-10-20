package model

type Message struct {
	ID      uint   `gorm:"primaryKey;type:bigint" json:"id"`
	Message string `json:"message"`
	Readed  bool   `json:"readed"`
}

func (Message) TableName() string {
	return "message"
}
