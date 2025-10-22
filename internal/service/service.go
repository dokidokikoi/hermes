package service

import "izumi/db"

type Iservice interface {
	Game() IGame
	Character() ICharacter
	Person() IPerson
	Policy() IPolicy
	Library() ILibrary
}

type service struct {
	store db.IStore
}

func (srv *service) Game() IGame {
	return NewGame(srv.store)
}

func (srv *service) Character() ICharacter {
	return NewCharacter(srv.store)
}

func (srv *service) Person() IPerson {
	return NewPerson(srv.store)
}

func (srv *service) Policy() IPolicy {
	return NewPolicy(srv.store)
}

func (srv *service) Library() ILibrary {
	return NewLibrary(srv.store)
}

func NewSrv(store db.IStore) *service {
	return &service{store: store}
}
