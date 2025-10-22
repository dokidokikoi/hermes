package model

type Message struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Message string `json:"message"`
	Readed  bool   `json:"readed"`
}

func (Message) TableName() string {
	return "message"
}
