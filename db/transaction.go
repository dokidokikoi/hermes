package db

import (
	"izumi/db/postgres"
)

// ITransaction 定义事务接口
type ITransaction interface {
	Begin() IStore
	Rollback() error
	Commit() error
}

// transaction 实现 ITransaction
type transaction struct {
	pg *postgres.Store
}

func newTransaction(pg *postgres.Store) *transaction {
	return &transaction{pg: pg}
}

func (t *transaction) Begin() IStore {
	return &Store{pg: t.pg.Begin()}
}

func (t *transaction) Rollback() error {
	return t.pg.Rollback()
}

func (t *transaction) Commit() error {
	return t.pg.Commit()
}