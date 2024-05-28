package service

import "hermes/db"

type Iservice interface {
	Game() IGame
}

type service struct {
	store db.IStore
}

func (srv *service) Game() IGame {
	return NewGame(srv.store)
}

func NewSrv(store db.IStore) *service {
	return &service{store: store}
}
