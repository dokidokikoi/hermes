package db

import (
	"hermes/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type IPolicy interface {
	base.BasicCURD[model.Policy]
}
