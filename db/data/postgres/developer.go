package postgres

import (
	"hermes/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type Developers struct {
	base.PgModel[model.Developer]
}

func NewDevelopers(db *Store) *Developers {
	return &Developers{PgModel: base.PgModel[model.Developer]{DB: db.DB}}
}

func init() {
	migrateTables = append(migrateTables, &model.Developer{})
}
