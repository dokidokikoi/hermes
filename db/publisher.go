package db

import (
	"hermes/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type IPublisher interface {
	base.BasicCURD[model.Publisher]
}
