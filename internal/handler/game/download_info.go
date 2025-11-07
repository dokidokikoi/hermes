package game

import (
	"context"
	"izumi/db/data"
	systemtask "izumi/internal/system_task"
	"izumi/model"
	"time"

	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/dokidokikoi/go-common/middleware"
)

type DownloadInfoReq struct {
	GameId uint `json:"game_id"`
}

func (h Handler) DownloadInfo(ctx context.Context, input *DownloadInfoReq, op *middleware.PreHandleOptions) (any, error) {
	err := h.srv.Game().DownloadInfo(ctx, input.GameId, time.Now())
	if err != nil {
		op.SetMsg(err.Error())
		return nil, err
	}
	return nil, nil
}

func (h Handler) DownloadAllInfo(ctx context.Context, input *struct{}, op *middleware.PreHandleOptions) (any, error) {
	gs, err := data.GetDataFactory().Game().List(ctx, &model.Game{}, &meta.ListOption{GetOption: meta.GetOption{Select: []string{"ID"}}})
	if err != nil {
		return nil, err
	}
	t := &model.SystemTask{
		Param:     model.SystemTaskParam{Amount: len(gs)},
		Type:      model.SystemTaskTypeDownload,
		State:     model.SystemTaskStateRunning,
		CreatedAt: time.Now(),
	}
	err = data.GetDataFactory().SystemTask().Create(ctx, t, nil)
	if err != nil {
		return nil, err
	}
	systemtask.DownloadTask(gs, t)

	return nil, nil
}
