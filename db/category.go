package db

import (
	"hermes/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type ICategory interface {
	base.BasicCURD[model.Category]
}
