package systemtask

import (
	"context"
	"izumi/constant"
	"izumi/db/data"
	"izumi/internal/service"
	"izumi/model"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/dokidokikoi/go-common/notice"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func init() {
	ts, err := data.GetDataFactory().SystemTask().List(context.Background(), &model.SystemTask{
		Type:  model.SystemTaskTypeDownload,
		State: model.SystemTaskStateRunning,
	}, &meta.ListOption{Order: "id desc"})
	if err != nil {
		zaplog.L().Error("system download task error", zap.Error(err))
		return
	}
	if len(ts) == 0 {
		return
	}
	err = data.GetDataFactory().SystemTask().UpdateByWhere(
		context.Background(),
		&meta.WhereNode{
			Conditions: []*meta.Condition{
				{
					Field:    "id",
					Operator: meta.NOTEQUAL,
					Value:    ts[0].ID,
				},
			},
		},
		&model.SystemTask{
			State: model.SystemTaskStateCanceled,
		},
		nil,
	)
	if err != nil {
		zaplog.L().Error("system update error", zap.Error(err))
	}

	gs, err := data.GetDataFactory().Game().List(context.Background(), &model.Game{}, &meta.ListOption{GetOption: meta.GetOption{Select: []string{"ID"}}})
	if err != nil {
		zaplog.L().Error("list game error", zap.Error(err))
		return
	}
	DownloadTask(gs, ts[0])
}

func DownloadTask(gs []*model.Game, t *model.SystemTask) {
	go func() {
		defer func() {
			err := data.GetDataFactory().SystemTask().Update(context.Background(), &model.SystemTask{
				ID:    t.ID,
				State: model.SystemTaskStateDone,
			}, nil)
			if err != nil {
				zaplog.L().Error("system update error", zap.Error(err))
			}
		}()
		srv := service.NewGame(data.GetDataFactory())
		for i, g := range gs {
			err := srv.DownloadInfo(context.Background(), g.ID, t.CreatedAt)
			if err != nil {
				zaplog.L().Error("downloadInfo", zap.Error(err))
			}
			notice.HubIns.SendBroadcast(constant.TOPIC_INFO_FILE, notice.NoticeResponse{
				Rid:     uuid.NewString(),
				Event:   constant.EVENT_INFO_FILE_DOWNLOAD,
				Success: err == nil,
				Message: func() string {
					if err != nil {
						return err.Error()
					}
					return ""
				}(),
				Data: map[string]any{
					"task_id":  t.ID,
					"game_id":  g.ID,
					"proccess": i + 1,
				},
			})
		}
	}()

}
