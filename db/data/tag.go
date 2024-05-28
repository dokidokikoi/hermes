package data

import (
	"hermes/db"
	"hermes/db/data/postgres"
)

var _ db.ITag = (*tag)(nil)

type tag struct {
	postgres.Tags
}

func newTag(d *data) *tag {
	return &tag{Tags: *postgres.NewTags(d.pg)}
}
