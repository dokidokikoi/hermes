package db

import (
	"hermes/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type IDeveloper interface {
	base.BasicCURD[model.Developer]
}
