package postgres

import (
	"hermes/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type Tags struct {
	base.PgModel[model.Tag]
}

func NewTags(db *Store) *Tags {
	return &Tags{PgModel: base.PgModel[model.Tag]{DB: db.DB.Debug()}}
}

func init() {
	migrateTables = append(migrateTables, &model.Tag{})
}
