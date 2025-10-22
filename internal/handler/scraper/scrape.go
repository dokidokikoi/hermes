package scraper

import (
	"context"
	"encoding/json"
	"izumi/constant"
	"izumi/db/data"
	"izumi/internal/handler"
	"izumi/model"
	"izumi/scraper/event"
	"time"

	"github.com/dokidokikoi/go-common/notice"

	"github.com/dokidokikoi/go-common/gopool"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (h Handler) Scrape(ctx context.Context, input *handler.ScraperDetailReq) (string, error) {
	if input.RequestID == "" {
		input.RequestID = uuid.New().String()
	}
	requestID := input.RequestID
	for _, req := range input.Objs {
		_, err := data.GetDataFactory().Task().Get(ctx, &model.Task{RequestID: requestID, Param: req.Url}, nil)
		if err == nil {
			continue
		}
		s, ok := event.GameScraperMap[req.Name]
		if !ok {
			continue
		}
		gopool.CtxGo(ctx, func() {
			var err error
			defer func() {
				msg := "success"
				if err != nil {
					msg = "failed"
				}
				notice.HubIns.SendBroadcast(constant.TOPIC_SCRAPER, notice.NoticeResponse{
					Rid:     requestID,
					Event:   constant.EVENT_SCRAPER_DETAIL,
					Message: msg,
					Data: map[string]any{
						"name": s.GetName(),
					},
				})
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
			item, err := s.GetItem(req.Url)
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
	return requestID, nil
}
