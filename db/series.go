package db

import (
	"izumi/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type ISeries interface {
	base.BasicCURD[model.Series]
}
