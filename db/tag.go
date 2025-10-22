package db

import (
	"izumi/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type ITag interface {
	base.BasicCURD[model.Tag]
}
