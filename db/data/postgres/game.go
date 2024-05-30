package postgres

import (
	"hermes/model"

	"github.com/dokidokikoi/go-common/db/base"
)

type Games struct {
	base.PgModel[model.Game]
}

func NewGames(db *Store) *Games {
	return &Games{PgModel: base.PgModel[model.Game]{DB: db.DB}}
}

type GameSeriess struct {
	base.PgModel[model.GameSeries]
}

func NewGameSeriess(db *Store) *GameSeriess {
	return &GameSeriess{PgModel: base.PgModel[model.GameSeries]{DB: db.DB}}
}

type GameCharacters struct {
	base.PgModel[model.GameCharacter]
}

func NewGameCharacters(db *Store) *GameCharacters {
	return &GameCharacters{PgModel: base.PgModel[model.GameCharacter]{DB: db.DB}}
}

type GameTags struct {
	base.PgModel[model.GameTag]
}

func NewGameTags(db *Store) *GameTags {
	return &GameTags{PgModel: base.PgModel[model.GameTag]{DB: db.DB}}
}

func init() {
	migrateTables = append(migrateTables, &model.Game{}, &model.GameCharacter{}, &model.GameSeries{}, &model.GameTag{})
}

type GameStaffs struct {
	base.PgModel[model.GameStaff]
}

func NewGameStaffs(db *Store) *GameStaffs {
	return &GameStaffs{PgModel: base.PgModel[model.GameStaff]{DB: db.DB}}
}

type GameInstances struct {
	base.PgModel[model.GameInstance]
}

func NewGameInstances(db *Store) *GameInstances {
	return &GameInstances{PgModel: base.PgModel[model.GameInstance]{DB: db.DB}}
}

func init() {
	migrateTables = append(migrateTables, &model.Game{}, &model.GameCharacter{}, &model.GameSeries{}, &model.GameTag{}, &model.GameStaff{}, &model.GameInstance{})
}
