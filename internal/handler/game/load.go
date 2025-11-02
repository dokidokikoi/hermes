package game

import (
	"context"
	"izumi/db/data"
	systemtask "izumi/internal/system_task"
	"izumi/model"
	"os"

	"github.com/dokidokikoi/go-common/middleware"
)

type LoadInfoResponse struct {
	Rid string `json:"rid"`
}

func (h Handler) LoadInfo(ctx context.Context, req *struct{}, op *middleware.PreHandleOptions) (any, error) {
	p, err := data.GetDataFactory().Policy().Get(ctx, &model.Policy{Key: model.SystemPolicy{}.Key()}, nil)
	if err != nil {
		return nil, err
	}
	sp, err := model.Parse[model.SystemPolicy](p.Policy)
	if err != nil {
		return nil, err
	}
	infos, err := h.srv.Library().Ls(ctx, sp.GameLibrary, false)
	if err != nil {
		if os.IsNotExist(err) {
			op.SetMsg("game library not exist")
		}
		return nil, err
	}
	t := &model.SystemTask{
		Amount: len(infos),
		Type:   model.SystemTaskTypeLoad,
		State:  model.SystemTaskStateRunning,
	}
	err = data.GetDataFactory().SystemTask().Create(ctx, t, nil)
	if err != nil {
		return nil, err
	}
	systemtask.LoadTask(infos, t)

	return nil, nil
}
