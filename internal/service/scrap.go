package service

import (
	"context"
	"encoding/json"
	"izumi/constant"
	"izumi/db"
	"izumi/db/data"
	"izumi/model"
	"izumi/scraper/event"
	"izumi/utils"
	"sync"
	"time"

	"github.com/dokidokikoi/go-common/gopool"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/dokidokikoi/go-common/notice"
	"go.uber.org/zap"
)

type IScrap interface {
	Scrap(ctx context.Context, requestID string, objs []model.ScrapObj, afterHook ...func(scraperName string, success bool)) (*sync.WaitGroup, error)
}

var _ IScrap = (*scrap)(nil)

type scrap struct {
	store db.IStore
}

func (s *scrap) Scrap(ctx context.Context, requestID string, objs []model.ScrapObj, afterHook ...func(scraperName string, success bool)) (*sync.WaitGroup, error) {
	wait := &sync.WaitGroup{}
	for _, req := range objs {
		_, err := data.GetDataFactory().Task().Get(ctx, &model.Task{RequestID: requestID, Param: req.Url}, nil)
		if err == nil {
			continue
		}
		s, ok := event.GameScraperMap[req.Name]
		if !ok {
			continue
		}
		wait.Add(1)
		gopool.CtxGo(ctx, func() {
			var err error
			defer func() {
				wait.Done()
				success := true
				if err != nil {
					success = false
				}
				notice.HubIns.SendBroadcast("", notice.NoticeResponse{
					Rid:     requestID,
					Event:   utils.ConcatEvent(constant.TOPIC_SCRAPER, constant.EVENT_SCRAPER_DETAIL),
					Success: success,
					Data: map[string]any{
						"name": s.GetName(),
					},
				})
				if len(afterHook) > 0 {
					afterHook[0](s.GetName(), success)
				}
			}()

			task := &model.Task{
				RequestID:   requestID,
				ScraperName: req.Name,
				Type:        model.TaskTypeDetail,
				Param:       req.Url,
				Status:      model.TaskStatusWait,
				StartAt:     time.Now(),
			}
			err = data.GetDataFactory().Task().Create(ctx, task, nil)
			if err != nil {
				zaplog.L().Error("创建任务失败", zap.Error(err))
				return
			}
			defer func() {
				err := data.GetDataFactory().Task().Update(ctx, task, nil)
				if err != nil {
					zaplog.L().Error("更新任务失败", zap.Error(err))
					return
				}
			}()
			item, err := s.GetGameItem(req.Url)
			task.EndAt = time.Now()
			task.Duration = int64(task.EndAt.Sub(task.StartAt).Seconds())
			if err != nil {
				task.Status = model.TaskStatusFailed
				zaplog.L().Error("获取详情失败", zap.String("scraper name", req.Name), zap.String("url", req.Url), zap.Error(err))
				return
			}

			data, err := json.Marshal(item)
			if err != nil {
				task.Status = model.TaskStatusFailed
				zaplog.L().Error("刮削失败", zap.Error(err))
				return
			}
			task.Result = string(data)
			task.Status = model.TaskStatusSucceed
		})
	}
	return wait, nil
}

func NewScrap(store db.IStore) *scrap {
	return &scrap{store: store}
}
