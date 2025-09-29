package db

type ITransaction interface {
	Begin() IStore
	Rollback() error
	Commit() error
}
