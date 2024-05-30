package data

import (
	"hermes/db"
	"hermes/db/data/postgres"
)

var _ db.IGame = (*game)(nil)

type game struct {
	postgres.Games
}

func newGame(d *data) *game {
	return &game{Games: *postgres.NewGames(d.pg)}
}

type gameTag struct {
	postgres.GameTags
}

func newGameTag(d *data) *gameTag {
	return &gameTag{GameTags: *postgres.NewGameTags(d.pg)}
}

type gameCharacter struct {
	postgres.GameCharacters
}

func newGameCharacter(d *data) *gameCharacter {
	return &gameCharacter{GameCharacters: *postgres.NewGameCharacters(d.pg)}
}

type gameSeries struct {
	postgres.GameSeriess
}

func newGameSeries(d *data) *gameSeries {
	return &gameSeries{GameSeriess: *postgres.NewGameSeriess(d.pg)}
}

type gameStaffs struct {
	postgres.GameStaffs
}

func newGameStaff(d *data) *gameStaffs {
	return &gameStaffs{GameStaffs: *postgres.NewGameStaffs(d.pg)}
}

type gameInstances struct {
	postgres.GameInstances
}

func newGameInstance(d *data) *gameInstances {
	return &gameInstances{GameInstances: *postgres.NewGameInstances(d.pg)}
}
