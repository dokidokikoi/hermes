package service

import "hermes/db"

type Iservice interface {
	Game() IGame
	Character() ICharacter
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

func NewSrv(store db.IStore) *service {
	return &service{store: store}
}
