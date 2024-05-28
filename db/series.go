package db

import (
	"hermes/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type ISeries interface {
	base.BasicCURD[model.Series]
}
