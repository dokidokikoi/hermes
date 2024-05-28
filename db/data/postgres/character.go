package postgres

import (
	"hermes/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type Characters struct {
	base.PgModel[model.Character]
}

func NewCharacters(db *Store) *Characters {
	return &Characters{PgModel: base.PgModel[model.Character]{DB: db.DB}}
}

type CharacterTags struct {
	base.PgModel[model.CharacterTag]
}

func NewCharacterTags(db *Store) *CharacterTags {
	return &CharacterTags{PgModel: base.PgModel[model.CharacterTag]{DB: db.DB}}
}

func init() {
	migrateTables = append(migrateTables, &model.Character{}, &model.CharacterTag{})
}
