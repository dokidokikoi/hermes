package data

import (
	"izumi/db"
	"izumi/db/data/postgres"
)

var _ db.IPolicy = (*policy)(nil)

type policy struct {
	postgres.Policies
}

func newPolicy(d *data) *policy {
	return &policy{Policies: *postgres.NewPolicies(d.pg)}
}
