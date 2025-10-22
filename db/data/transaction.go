package data

import (
	"izumi/db"
	"izumi/db/data/postgres"
)

var _ db.ITransaction = (*transaction)(nil)

type transaction struct {
	pg *postgres.Store
}

func (d *transaction) Begin() db.IStore {
	return &data{pg: d.pg.TransactionBegin()}
}

func (d *transaction) Rollback() error {
	return d.pg.TransactionRollback()
}

func (d *transaction) Commit() error {
	return d.pg.TransactionCommit()
}

func newTransaction(d *data) *transaction {
	return &transaction{pg: d.pg}
}
