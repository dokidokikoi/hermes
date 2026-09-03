package postgres

import (
	"izumi/config"

	"github.com/dokidokikoi/go-common/db"
	"gorm.io/gorm"
)

// Store 封装数据库连接
type Store struct {
	DB *gorm.DB
}

// NewStore 创建数据库 Store
func NewStore() *Store {
	cfg := config.GetConfig()
	var dbIns *gorm.DB
	var err error

	if cfg.PGConfig.Database == "" {
		// 使用 SQLite
		dbIns, err = db.NewSqlite(cfg.SqliteConfig.Database)
	} else {
		// 使用 PostgreSQL
		dbIns, err = db.NewPostgresql(
			cfg.PGConfig.Username,
			cfg.PGConfig.Database,
			db.WithHost(cfg.PGConfig.Host),
			db.WithPort(cfg.PGConfig.Port),
			db.WithPassword(cfg.PGConfig.Password),
		)
	}

	if err != nil {
		panic(err)
	}

	dbIns = dbIns.Session(&gorm.Session{FullSaveAssociations: true})

	// 执行迁移
	if err := migrateDatabase(dbIns); err != nil {
		panic(err)
	}

	return &Store{DB: dbIns.Debug()}
}

// Begin 开始事务
func (s *Store) Begin() *Store {
	return &Store{DB: s.DB.Begin()}
}

// Rollback 回滚事务
func (s *Store) Rollback() error {
	return s.DB.Rollback().Error
}

// Commit 提交事务
func (s *Store) Commit() error {
	return s.DB.Commit().Error
}

// Close 关闭数据库连接
func (s *Store) Close() error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
