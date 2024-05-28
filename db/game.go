package db

import (
	"hermes/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type IGame interface {
	base.BasicCURD[model.Game]
}

type IGameCharacter interface {
	base.BasicCURD[model.GameCharacter]
}

type IGameSeries interface {
	base.BasicCURD[model.GameSeries]
}

type IGameTag interface {
	base.BasicCURD[model.GameTag]
}
