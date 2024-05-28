package data

import (
	"hermes/db"
	"hermes/db/data/postgres"
)

var _ db.IPublisher = (*publisher)(nil)

type publisher struct {
	postgres.Publishers
}

func newPublisher(d *data) *publisher {
	return &publisher{Publishers: *postgres.NewPublishers(d.pg)}
}
