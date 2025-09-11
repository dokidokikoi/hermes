package game

import (
	"context"
	"hermes/db/data"
	"hermes/internal/handler"
	"hermes/model"
)

type GetVerRequest struct {
	ID uint `form:"id" query:"id" json:"id"`
}

func (h Handler) Get(ctx context.Context, req *GetVerRequest) (*handler.GameVo, error) {
	// veriosn
	gVo, err := h.srv.Game().GetVOByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return gVo, nil
}

type GetInsRequest struct {
	ID uint `form:"id" query:"id" json:"id"`
}

func (h Handler) GetIns(ctx context.Context, req *GetInsRequest) ([]*model.GameInstance, error) {
	gIns, err := data.GetDataFactory().GameInstance().List(ctx, &model.GameInstance{GameID: req.ID}, nil)
	if err != nil {
		return nil, err
	}

	return gIns, nil
}
