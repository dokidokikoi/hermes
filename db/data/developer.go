package data

import (
	"hermes/db"
	"hermes/db/data/postgres"
)

var _ db.IDeveloper = (*developer)(nil)

type developer struct {
	postgres.Developers
}

func newDeveloper(d *data) *developer {
	return &developer{Developers: *postgres.NewDevelopers(d.pg)}
}
