package postgres

import (
	"hermes/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type Publishers struct {
	base.PgModel[model.Publisher]
}

func NewPublishers(db *Store) *Publishers {
	return &Publishers{PgModel: base.PgModel[model.Publisher]{DB: db.DB}}
}

func init() {
	migrateTables = append(migrateTables, &model.Publisher{})
}
