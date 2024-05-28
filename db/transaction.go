package db

type ITransaction interface {
	Begin() IStore
	Rollback()
	Commit()
}
