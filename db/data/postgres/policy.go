package postgres

import (
	"hermes/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type Policies struct {
	base.PgModel[model.Policy]
}

func NewPolicies(db *Store) *Policies {
	return &Policies{PgModel: base.PgModel[model.Policy]{DB: db.DB}}
}

func init() {
	migrateTables = append(migrateTables, &model.Policy{})
}
