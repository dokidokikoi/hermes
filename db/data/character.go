package data

import (
	"hermes/db"
	"hermes/db/data/postgres"
)

var _ db.ICharacter = (*character)(nil)

type character struct {
	postgres.Characters
}

func newCharacter(d *data) *character {
	return &character{Characters: *postgres.NewCharacters(d.pg)}
}

var _ db.ICharacter = (*character)(nil)

type characterTag struct {
	postgres.CharacterTags
}

func newCharacterTag(d *data) *characterTag {
	return &characterTag{CharacterTags: *postgres.NewCharacterTags(d.pg)}
}
