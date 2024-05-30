package data

import (
	"hermes/db"
	"hermes/db/data/postgres"
)

var _ db.IPerson = (*person)(nil)

type person struct {
	postgres.Persons
}

func newPerson(d *data) *person {
	return &person{Persons: *postgres.NewPersons(d.pg)}
}
