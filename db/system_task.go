package db

import (
	"izumi/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type ISystemTask interface {
	base.BasicCURD[model.SystemTask]
}
