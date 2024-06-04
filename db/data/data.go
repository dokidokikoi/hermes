package data

import (
	"hermes/config"
	"hermes/db"
	"hermes/db/data/postgres"
	"sync"
)

var _ db.IStore = (*data)(nil)

var (
	dataIns *data
	once    sync.Once
)

type data struct {
	pg *postgres.Store
}

func (d *data) Category() db.ICategory {
	return newCategory(d)
}
func (d *data) Character() db.ICharacter {
	return newCharacter(d)
}
func (d *data) CharacterTag() db.ICharacterTag {
	return newCharacterTag(d)
}
func (d *data) Developer() db.IDeveloper {
	return newDeveloper(d)
}
func (d *data) Publisher() db.IPublisher {
	return newPublisher(d)
}
func (d *data) Series() db.ISeries {
	return newSeries(d)
}
func (d *data) Tag() db.ITag {
	return newTag(d)
}
func (d *data) Game() db.IGame {
	return newGame(d)
}
func (d *data) GameCharacter() db.IGameCharacter {
	return newGameCharacter(d)
}
func (d *data) GameSeries() db.IGameSeries {
	return newGameSeries(d)
}
func (d *data) GameTag() db.IGameTag {
	return newGameTag(d)
}
func (d *data) Person() db.IPerson {
	return newPerson(d)
}
func (d *data) GameStaff() db.IGameStaff {
	return newGameStaff(d)
}
func (d *data) GameInstance() db.IGameInstance {
	return newGameInstance(d)
}
func (d *data) Task() db.ITask {
	return newTask(d)
}
func (d *data) RefGameInstance() db.IRefGameInstance {
	return newRefGameInstance(d)
}
func (d *data) Transaction() db.ITransaction {
	return newTransaction(d)
}

func GetDataFactory() *data {
	once.Do(func() {
		dataIns = &data{
			pg: postgres.NewPostgresStore(config.GetConfig().PGConfig),
		}
	})

	return &data{pg: dataIns.pg}
}

func Close() error {
	if err := dataIns.pg.Close(); err != nil {
		return err
	}
	return nil
}
