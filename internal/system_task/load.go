package systemtask

import (
	"context"
	"encoding/json"
	"izumi/constant"
	"izumi/db/data"
	"izumi/internal/handler"
	"izumi/internal/service"
	"izumi/model"
	"os"
	"path/filepath"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/dokidokikoi/go-common/notice"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func init() {
	ts, err := data.GetDataFactory().SystemTask().List(context.Background(), &model.SystemTask{
		Type:  model.SystemTaskTypeLoad,
		State: model.SystemTaskStateRunning,
	}, &meta.ListOption{Order: "id desc"})
	if err != nil {
		zaplog.L().Error("system load task error", zap.Error(err))
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
		return
	}

	p, err := data.GetDataFactory().Policy().Get(context.Background(), &model.Policy{Key: model.SystemPolicy{}.Key()}, nil)
	if err != nil {
		zaplog.L().Error("get system policy error", zap.Error(err))
		return
	}
	sp, err := model.Parse[model.SystemPolicy](p.Policy)
	if err != nil {
		zaplog.L().Error("parse system policy error", zap.Error(err))
		return
	}

	srv := service.NewSrv(data.GetDataFactory())
	infos, err := srv.Library().Ls(context.Background(), sp.GameLibrary)
	if err != nil {
		zaplog.L().Error("get game library error", zap.Error(err))
	}
	LoadTask(infos, ts[0])
}

func LoadTask(infos []service.PathInfo, t *model.SystemTask) {
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
		for i, info := range infos {
			err := load(context.Background(), info)
			if err != nil {
				zaplog.L().Error("load info file failed", zap.String("path", info.Path), zap.Error(err))
			}
			notice.HubIns.SendBroadcast(constant.TOPIC_INFO_FILE, notice.NoticeResponse{
				Rid:     uuid.NewString(),
				Event:   constant.EVENT_INFO_FILE_LOAD,
				Success: err == nil,
				Message: func() string {
					if err != nil {
						return err.Error()
					}
					return ""
				}(),
				Data: map[string]any{
					"task_id":  t.ID,
					"path":     info.Path,
					"proccess": i + 1,
				},
			})
		}
	}()
}

func load(ctx context.Context, info service.PathInfo) error {
	if !info.IsDir {
		return nil
	}
	path := filepath.Join(info.Path, "info.json")
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	gVo := &handler.GameVo{}
	err = json.NewDecoder(f).Decode(gVo)
	if err != nil {
		return err
	}

	return service.NewGame(data.GetDataFactory()).Load(ctx, gVo, info.Path)
}
