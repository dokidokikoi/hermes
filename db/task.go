package db

import (
	"izumi/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type ITask interface {
	base.BasicCURD[model.Task]
}
