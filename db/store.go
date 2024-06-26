package db

type IStore interface {
	Category() ICategory
	Character() ICharacter
	CharacterTag() ICharacterTag
	Developer() IDeveloper
	Publisher() IPublisher
	Series() ISeries
	Tag() ITag
	Game() IGame
	GameCharacter() IGameCharacter
	GameSeries() IGameSeries
	GameTag() IGameTag
	Transaction() ITransaction
	Person() IPerson
	PersonTag() IPersonTag
	GameStaff() IGameStaff
	GameInstance() IGameInstance
	Task() ITask
	RefGameInstance() IRefGameInstance
	Policy() IPolicy
}
