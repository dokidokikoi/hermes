package data

import (
	"izumi/db"
	"izumi/db/data/postgres"
)

var _ db.ISystemTask = (*systemTask)(nil)

type systemTask struct {
	postgres.SystemTasks
}

func newSystemTask(d *data) *systemTask {
	return &systemTask{SystemTasks: *postgres.NewSystemTasks(d.pg)}
}
