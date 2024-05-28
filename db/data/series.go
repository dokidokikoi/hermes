package data

import (
	"hermes/db"
	"hermes/db/data/postgres"
)

var _ db.ISeries = (*series)(nil)

type series struct {
	postgres.Seriess
}

func newSeries(d *data) *series {
	return &series{Seriess: *postgres.NewSeriess(d.pg)}
}
