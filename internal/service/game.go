package service

import (
	"context"
	"hermes/db"
	"hermes/model"
)

type IGame interface {
	CreateL(ctx context.Context, g *model.Game) error
	UpdateL(ctx context.Context, g *model.Game) error
}

type game struct {
	store db.IStore
}

func (gsrv *game) CreateL(ctx context.Context, g *model.Game) error {
	return gsrv.store.Game().Create(ctx, g, nil)
}

func (gsrv *game) UpdateL(ctx context.Context, g *model.Game) error {
	tx := gsrv.store.Transaction().Begin()
	err := tx.GameTag().Delete(ctx, &model.GameTag{GameID: g.ID}, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}
	err = tx.GameSeries().Delete(ctx, &model.GameSeries{GameID: g.ID}, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}
	err = tx.GameCharacter().Delete(ctx, &model.GameCharacter{GameID: g.ID}, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}

	err = tx.Game().Update(ctx, g, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}
	return nil
}

func NewGame(store db.IStore) *game {
	return &game{store: store}
}
