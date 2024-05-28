package postgres

import (
	"hermes/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type Seriess struct {
	base.PgModel[model.Series]
}

func NewSeriess(db *Store) *Seriess {
	return &Seriess{PgModel: base.PgModel[model.Series]{DB: db.DB}}
}

func init() {
	migrateTables = append(migrateTables, &model.Series{})
}
