package db

import (
	"github.com/dokidokikoi/go-common/db/base"
	"gorm.io/gorm"
)

// IRepo 泛型接口，保持类型安全的 CRUD 操作
type IRepo[T base.BaseModel] interface {
	base.BasicCURD[T]
}

// PgRepo 泛型实现，直接嵌入 base.PgModel
type PgRepo[T base.BaseModel] struct {
	base.PgModel[T]
}

// NewRepo 创建泛型 Repository
func NewRepo[T base.BaseModel](db *gorm.DB) IRepo[T] {
	return &PgRepo[T]{PgModel: base.PgModel[T]{DB: db}}
}