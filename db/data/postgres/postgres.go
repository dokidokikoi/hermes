package postgres

import (
	"github.com/dokidokikoi/go-common/config"
	"github.com/dokidokikoi/go-common/db"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var migrateTables []any
var cleanTables []any

type Store struct {
	DB *gorm.DB
}

func (t *Store) TransactionBegin() *Store {
	db := t.DB.Begin()
	return &Store{
		DB: db,
	}
}

func (t *Store) TransactionRollback() {
	t.DB.Rollback()
}

func (t *Store) TransactionCommit() {
	t.DB.Commit()
}

func (t *Store) Close() error {
	sqlDB, err := t.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func NewPostgresStore(c config.PGConfig) *Store {
	dbIns, err := db.NewPostgresql(c.Username, c.Database, db.WithHost(c.Host), db.WithPort(c.Port), db.WithPassword(c.Password))
	if err != nil {
		panic(err)
	}
	dbIns = dbIns.Session(&gorm.Session{FullSaveAssociations: true})
	pgFactory := &Store{dbIns}

	// cleanDatabase(dbIns)
	// 自动化迁移
	if err := migrateDatabase(dbIns); err != nil {
		panic(err)
	}
	return pgFactory
}

func migrateDatabase(db *gorm.DB) error {
	for _, t := range migrateTables {
		if err := db.AutoMigrate(t); err != nil {
			return errors.Wrap(err, "migrate model failed")
		}
	}
	// if err := db.AutoMigrate(&model.Category{}); err != nil {
	// 	return errors.Wrap(err, "migrate model failed")
	// }
	// if err := db.AutoMigrate(&model.Tag{}); err != nil {
	// 	return errors.Wrap(err, "migrate model failed")
	// }
	// if err := db.AutoMigrate(&model.Series{}); err != nil {
	// 	return errors.Wrap(err, "migrate model failed")
	// }
	// if err := db.AutoMigrate(&model.Publisher{}); err != nil {
	// 	return errors.Wrap(err, "migrate model failed")
	// }
	// if err := db.AutoMigrate(&model.Developer{}); err != nil {
	// 	return errors.Wrap(err, "migrate model failed")
	// }
	// if err := db.AutoMigrate(&model.Character{}); err != nil {
	// 	return errors.Wrap(err, "migrate model failed")
	// }
	// if err := db.AutoMigrate(&model.Game{}); err != nil {
	// 	return errors.Wrap(err, "migrate model failed")
	// }
	return nil
}

func cleanDatabase(db *gorm.DB) error {
	for _, t := range migrateTables {
		if err := db.Migrator().DropTable(t); err != nil {
			return errors.Wrap(err, "drop model failed")
		}
	}
	return nil
}
