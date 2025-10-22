package data

import (
	"izumi/db"
	"izumi/db/data/postgres"
)

var _ db.ITask = (*task)(nil)

type task struct {
	postgres.Tasks
}

func newTask(d *data) *task {
	return &task{Tasks: *postgres.NewTasks(d.pg)}
}
