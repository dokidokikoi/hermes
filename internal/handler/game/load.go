package game

import (
	"context"
	"izumi/db/data"
	"izumi/internal/service"
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

	var infos []service.PathInfo
	for _, l := range sp.GameLibrary {
		is, err := h.srv.Library().Ls(ctx, l, false)
		if err != nil {
			if os.IsNotExist(err) {
				op.SetMsg("game library not exist")
			}
			return nil, err
		}
		infos = append(infos, is...)
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
