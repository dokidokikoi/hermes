package game

import (
	"context"
	"izumi/db"
	"izumi/internal/handler"
	"izumi/model"

	"github.com/dokidokikoi/go-common/middleware"
)

type GetVerRequest struct {
	ID uint `form:"id" query:"id" json:"id"`
}

func (h Handler) Get(ctx context.Context, req *GetVerRequest, op *middleware.PreHandleOptions) (*handler.GameVo, error) {
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

func (h Handler) GetIns(ctx context.Context, req *GetInsRequest, op *middleware.PreHandleOptions) ([]*model.GameInstance, error) {
	gIns, err := db.GetStore().GameInstance().List(ctx, &model.GameInstance{GameID: req.ID}, nil)
	if err != nil {
		return nil, err
	}

	return gIns, nil
}
