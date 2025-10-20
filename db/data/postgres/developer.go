package postgres

import (
	"hermes/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type Brands struct {
	base.PgModel[model.Brand]
}

func NewBrands(db *Store) *Brands {
	return &Brands{PgModel: base.PgModel[model.Brand]{DB: db.DB}}
}

func init() {
	migrateTables = append(migrateTables, &model.Brand{})
}
