package data

import (
	"izumi/config"
	"izumi/db"
	"izumi/db/data/postgres"
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
func (d *data) Brand() db.IBrand {
	return newBrand(d)
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
func (d *data) PersonTag() db.IPersonTag {
	return newPersonTag(d)
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

func (d *data) Policy() db.IPolicy {
	return newPolicy(d)
}

func (d *data) SystemTask() db.ISystemTask {
	return newSystemTask(d)
}

func (d *data) Transaction() db.ITransaction {
	return newTransaction(d)
}

func GetDataFactory() *data {
	once.Do(func() {
		var pg *postgres.Store
		if config.GetConfig().PGConfig.Database == "" {
			pg = postgres.NewSqliteStore(config.GetConfig().SqliteConfig)
		} else {
			pg = postgres.NewPostgresStore(config.GetConfig().PGConfig)
		}
		dataIns = &data{
			pg: pg,
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
