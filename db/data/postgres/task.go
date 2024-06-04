package postgres

import (
	"hermes/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type Tasks struct {
	base.PgModel[model.Task]
}

func NewTasks(db *Store) *Tasks {
	return &Tasks{PgModel: base.PgModel[model.Task]{DB: db.DB}}
}

type RefGameInstances struct {
	base.PgModel[model.RefGameInstance]
}

func NewRefGameInstances(db *Store) *RefGameInstances {
	return &RefGameInstances{PgModel: base.PgModel[model.RefGameInstance]{DB: db.DB}}
}

func init() {
	migrateTables = append(migrateTables, &model.Task{}, &model.RefGameInstance{})
}
