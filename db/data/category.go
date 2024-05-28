package data

import (
	"hermes/db"
	"hermes/db/data/postgres"
)

var _ db.ICategory = (*category)(nil)

type category struct {
	postgres.Categorys
}

func newCategory(d *data) *category {
	return &category{Categorys: *postgres.NewCategorys(d.pg)}
}
