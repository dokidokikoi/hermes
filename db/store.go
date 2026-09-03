package db

import (
	"izumi/db/postgres"
	"izumi/model"
	"sync"
)

// IStore 定义所有数据访问方法
type IStore interface {
	// Metadata domain
	Category() IRepo[model.Category]
	Tag() IRepo[model.Tag]
	Etag() IRepo[model.Etag]
	DecidedTag() IRepo[model.DecidedTag]
	Series() IRepo[model.Series]
	Brand() IRepo[model.Brand]

	// Character domain
	Character() IRepo[model.Character]
	CharacterTag() IRepo[model.CharacterTag]

	// Person domain
	Person() IRepo[model.Person]
	PersonTag() IRepo[model.PersonTag]

	// Game domain
	Game() IRepo[model.Game]
	GameCharacter() IRepo[model.GameCharacter]
	GameSeries() IRepo[model.GameSeries]
	GameTag() IRepo[model.GameTag]
	GameBrands() IRepo[model.GameBrands]
	GameStaff() IRepo[model.GameStaff]
	GameInstance() IRepo[model.GameInstance]

	// System domain
	Task() IRepo[model.Task]
	Policy() IRepo[model.Policy]
	SystemTask() IRepo[model.SystemTask]

	Transaction() ITransaction
}

// Store 实现 IStore
type Store struct {
	pg *postgres.Store
}

var (
	storeIns *Store
	once     sync.Once
)

// GetStore 获取单例 Store（修复原 bug：返回单例而非新实例）
func GetStore() *Store {
	once.Do(func() {
		storeIns = NewStore()
	})
	return storeIns
}

// NewStore 创建新 Store
func NewStore() *Store {
	return &Store{pg: postgres.NewStore()}
}

// Close 关闭数据库连接
func Close() error {
	if storeIns != nil {
		return storeIns.pg.Close()
	}
	return nil
}

// Metadata domain methods
func (s *Store) Category() IRepo[model.Category] {
	return NewRepo[model.Category](s.pg.DB)
}

func (s *Store) Tag() IRepo[model.Tag] {
	return NewRepo[model.Tag](s.pg.DB)
}

func (s *Store) Etag() IRepo[model.Etag] {
	return NewRepo[model.Etag](s.pg.DB)
}

func (s *Store) DecidedTag() IRepo[model.DecidedTag] {
	return NewRepo[model.DecidedTag](s.pg.DB)
}

func (s *Store) Series() IRepo[model.Series] {
	return NewRepo[model.Series](s.pg.DB)
}

func (s *Store) Brand() IRepo[model.Brand] {
	return NewRepo[model.Brand](s.pg.DB)
}

// Character domain methods
func (s *Store) Character() IRepo[model.Character] {
	return NewRepo[model.Character](s.pg.DB)
}

func (s *Store) CharacterTag() IRepo[model.CharacterTag] {
	return NewRepo[model.CharacterTag](s.pg.DB)
}

// Person domain methods
func (s *Store) Person() IRepo[model.Person] {
	return NewRepo[model.Person](s.pg.DB)
}

func (s *Store) PersonTag() IRepo[model.PersonTag] {
	return NewRepo[model.PersonTag](s.pg.DB)
}

// Game domain methods
func (s *Store) Game() IRepo[model.Game] {
	return NewRepo[model.Game](s.pg.DB)
}

func (s *Store) GameCharacter() IRepo[model.GameCharacter] {
	return NewRepo[model.GameCharacter](s.pg.DB)
}

func (s *Store) GameSeries() IRepo[model.GameSeries] {
	return NewRepo[model.GameSeries](s.pg.DB)
}

func (s *Store) GameTag() IRepo[model.GameTag] {
	return NewRepo[model.GameTag](s.pg.DB)
}

func (s *Store) GameBrands() IRepo[model.GameBrands] {
	return NewRepo[model.GameBrands](s.pg.DB)
}

func (s *Store) GameStaff() IRepo[model.GameStaff] {
	return NewRepo[model.GameStaff](s.pg.DB)
}

func (s *Store) GameInstance() IRepo[model.GameInstance] {
	return NewRepo[model.GameInstance](s.pg.DB)
}

// System domain methods
func (s *Store) Task() IRepo[model.Task] {
	return NewRepo[model.Task](s.pg.DB)
}

func (s *Store) Policy() IRepo[model.Policy] {
	return NewRepo[model.Policy](s.pg.DB)
}

func (s *Store) SystemTask() IRepo[model.SystemTask] {
	return NewRepo[model.SystemTask](s.pg.DB)
}

// Transaction method
func (s *Store) Transaction() ITransaction {
	return newTransaction(s.pg)
}
