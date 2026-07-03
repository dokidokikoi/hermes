package systemtask

import (
	"context"
	"fmt"
	"izumi/constant"
	"izumi/db"
	"izumi/internal/service"
	"izumi/model"
	"izumi/utils"
	"time"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/dokidokikoi/go-common/notice"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func StartDownload() {
	ts, err := db.GetStore().SystemTask().List(context.Background(), &model.SystemTask{
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
	err = db.GetStore().SystemTask().UpdateByWhere(
		context.Background(),
		&meta.WhereNode{
			Conditions: []*meta.Condition{
				{
					Field:    "id",
					Operator: meta.NOTEQUAL,
					Value:    ts[0].ID,
				},
			},
			Next: &meta.WhereNode{
				Conditions: []*meta.Condition{
					{
						Field:    "state",
						Operator: meta.EQUAL,
						Value:    model.SystemTaskStateRunning,
					},
				},
				Next: &meta.WhereNode{
					Conditions: []*meta.Condition{
						{
							Field:    "type",
							Operator: meta.EQUAL,
							Value:    model.SystemTaskTypeDownload,
						},
					},
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

	gs, err := db.GetStore().Game().List(context.Background(), &model.Game{}, &meta.ListOption{GetOption: meta.GetOption{Select: []string{"ID"}}})
	if err != nil {
		zaplog.L().Error("list game error", zap.Error(err))
		return
	}
	DownloadTask(gs, ts[0])
}

func DownloadTask(gs []*model.Game, t *model.SystemTask) {
	rid := uuid.NewString()
	go func() {
		defer func() {
			err := db.GetStore().SystemTask().Update(context.Background(), &model.SystemTask{
				ID:    t.ID,
				State: model.SystemTaskStateDone,
			}, nil)
			if err != nil {
				zaplog.L().Error("system update error", zap.Error(err))
			}
			notice.HubIns.SendBroadcast("", notice.NoticeResponse{
				Rid:     rid,
				Event:   utils.ConcatEvent(constant.TOPIC_INFO_FILE, constant.EVENT_INFO_FILE_DOWNLOAD),
				Success: true,
			})
		}()
		var (
			gid      uint
			proccess int
			total    = len(gs)
			message  string
		)
		ctxWithCancel, cancel := context.WithCancel(context.TODO())
		defer func() {
			cancel()
			notice.HubIns.SendBroadcast("", notice.NoticeResponse{
				Rid:     uuid.NewString(),
				Event:   utils.ConcatEvent(constant.TOPIC_SCRAPER, constant.EVENT_INFO_FILE_DOWNLOADING),
				Success: true,
				Data: map[string]any{
					"task_id":  t.ID,
					"proccess": total,
					"total":    total,
				},
			})
		}()
		go func() {
			ticker := time.NewTicker(time.Millisecond * 500)
			for {
				select {
				case <-ticker.C:
					notice.HubIns.SendBroadcast("", notice.NoticeResponse{
						Rid:     rid,
						Event:   utils.ConcatEvent(constant.TOPIC_INFO_FILE, constant.EVENT_INFO_FILE_DOWNLOADING),
						Success: message == "",
						Message: message,
						Data: map[string]any{
							"task_id":  t.ID,
							"game_id":  gid,
							"proccess": proccess,
							"total":    total,
						},
					})
				case <-ctxWithCancel.Done():
					return
				}
			}
		}()
		srv := service.NewGame(db.GetStore())
		for i, g := range gs {
			err := srv.DownloadInfo(context.Background(), g.ID, t.CreatedAt)
			if err != nil {
				zaplog.L().Error("downloadInfo", zap.Error(err))
				message = fmt.Sprintf("download '%s' error", g.Name)
			} else {
				message = ""
			}
			proccess = i + 1
			gid = g.ID
		}
	}()
}
