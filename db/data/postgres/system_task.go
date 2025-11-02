package postgres

import (
	"izumi/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type SystemTasks struct {
	base.PgModel[model.SystemTask]
}

func NewSystemTasks(db *Store) *SystemTasks {
	return &SystemTasks{PgModel: base.PgModel[model.SystemTask]{DB: db.DB}}
}

func init() {
	migrateTables = append(migrateTables, &model.SystemTask{})
}
