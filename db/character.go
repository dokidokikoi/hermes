package db

import (
	"izumi/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type ICharacter interface {
	base.BasicCURD[model.Character]
}

type ICharacterTag interface {
	base.BasicCURD[model.CharacterTag]
}
