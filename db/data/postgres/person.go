package postgres

import (
	"hermes/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type Persons struct {
	base.PgModel[model.Person]
}

func NewPersons(db *Store) *Persons {
	return &Persons{PgModel: base.PgModel[model.Person]{DB: db.DB}}
}

type PersonTags struct {
	base.PgModel[model.PersonTag]
}

func NewPersonTags(db *Store) *PersonTags {
	return &PersonTags{PgModel: base.PgModel[model.PersonTag]{DB: db.DB.Debug()}}
}

func init() {
	migrateTables = append(migrateTables, &model.Person{}, &model.PersonTag{})
}
