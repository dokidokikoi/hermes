package data

import (
	"izumi/db"
	"izumi/db/data/postgres"
)

var _ db.IBrand = (*brand)(nil)

type brand struct {
	postgres.Brands
}

func newBrand(d *data) *brand {
	return &brand{Brands: *postgres.NewBrands(d.pg)}
}
