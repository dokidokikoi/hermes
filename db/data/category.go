package data

import (
	"izumi/db"
	"izumi/db/data/postgres"
)

var _ db.ICategory = (*category)(nil)

type category struct {
	postgres.Categorys
}

func newCategory(d *data) *category {
	return &category{Categorys: *postgres.NewCategorys(d.pg)}
}
