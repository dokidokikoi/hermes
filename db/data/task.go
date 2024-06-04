package data

import (
	"hermes/db"
	"hermes/db/data/postgres"
)

var _ db.ITask = (*task)(nil)

type task struct {
	postgres.Tasks
}

func newTask(d *data) *task {
	return &task{Tasks: *postgres.NewTasks(d.pg)}
}

type refGameInstance struct {
	postgres.RefGameInstances
}

func newRefGameInstance(d *data) *refGameInstance {
	return &refGameInstance{RefGameInstances: *postgres.NewRefGameInstances(d.pg)}
}
