package db

import (
	"hermes/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type IPerson interface {
	base.BasicCURD[model.Person]
}
