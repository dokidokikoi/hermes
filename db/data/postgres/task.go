package postgres

import (
	"izumi/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type Tasks struct {
	base.PgModel[model.Task]
}

func NewTasks(db *Store) *Tasks {
	return &Tasks{PgModel: base.PgModel[model.Task]{DB: db.DB}}
}

func init() {
	migrateTables = append(migrateTables, &model.Task{})
}
