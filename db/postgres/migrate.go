package postgres

import (
	"izumi/model"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// migrateTables 集中定义所有需要迁移的表
var migrateTables = []any{
	// Metadata domain
	&model.Category{},
	&model.Tag{},
	&model.Series{},
	&model.Brand{},

	// Character domain
	&model.Character{},
	&model.CharacterTag{},

	// Person domain
	&model.Person{},
	&model.PersonTag{},

	// Game domain
	&model.Game{},
	&model.GameCharacter{},
	&model.GameSeries{},
	&model.GameTag{},
	&model.GameBrands{},
	&model.GameStaff{},
	&model.GameInstance{},

	// System domain
	&model.Task{},
	&model.Policy{},
	&model.SystemTask{},
}

// migrateDatabase 执行数据库迁移
func migrateDatabase(db *gorm.DB) error {
	for _, t := range migrateTables {
		if err := db.AutoMigrate(t); err != nil {
			return errors.Wrap(err, "migrate model failed")
		}
	}
	return nil
}

// cleanDatabase 清理数据库（用于开发/测试）
func cleanDatabase(db *gorm.DB) error {
	for _, t := range migrateTables {
		if err := db.Migrator().DropTable(t); err != nil {
			return errors.Wrap(err, "drop table failed")
		}
	}
	return nil
}