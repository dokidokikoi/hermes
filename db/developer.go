package db

import (
	"hermes/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type IBrand interface {
	base.BasicCURD[model.Brand]
}
