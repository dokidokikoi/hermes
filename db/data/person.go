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

var _ db.IPersonTag = (*personTag)(nil)

type personTag struct {
	postgres.PersonTags
}

func newPersonTag(d *data) *personTag {
	return &personTag{PersonTags: *postgres.NewPersonTags(d.pg)}
}
