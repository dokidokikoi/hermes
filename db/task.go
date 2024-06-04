package db

import (
	"hermes/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type ITask interface {
	base.BasicCURD[model.Task]
}

type IRefGameInstance interface {
	base.BasicCURD[model.RefGameInstance]
}
