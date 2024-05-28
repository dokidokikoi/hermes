package postgres

import (
	"hermes/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type Categorys struct {
	base.PgModel[model.Category]
}

func NewCategorys(db *Store) *Categorys {
	return &Categorys{PgModel: base.PgModel[model.Category]{DB: db.DB}}
}

func init() {
	migrateTables = append(migrateTables, &model.Category{})
}
