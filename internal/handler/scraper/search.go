package scraper

import (
	"context"
	"encoding/json"
	"izumi/constant"
	"izumi/db/data"
	"izumi/internal/handler"
	"izumi/model"
	"izumi/scraper"
	"izumi/scraper/event"
	"time"

	"github.com/dokidokikoi/go-common/middleware"
	"github.com/dokidokikoi/go-common/notice"

	comm_errs "github.com/dokidokikoi/go-common/errors"
	"github.com/dokidokikoi/go-common/gopool"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (h Handler) Search(ctx context.Context, input *handler.ScraperSearchReq, op *middleware.PreHandleOptions) (string, error) {
	if input.RequestID == "" {
		input.RequestID = uuid.New().String()
	}
	if input.Name == "all" {
		for _, s := range event.GameScraperMap {
			s := s
			gopool.CtxGo(ctx, func() {
				DoSearch(ctx, input.RequestID, *input, s)
			})
		}
	} else {
		if _, ok := event.GameScraperMap[input.Name]; !ok {
			op.SetErr(comm_errs.ApiErrValidation)
			return "", errors.New("scraper name not found")
		}
		gopool.CtxGo(ctx, func() {
			DoSearch(ctx, input.RequestID, *input, event.GameScraperMap[input.Name])
		})
	}
	return input.RequestID, nil
}

func DoSearch(ctx context.Context, requestID string, input handler.ScraperSearchReq, s scraper.IGameScraper) {
	var err error
	defer func() {
		msg := "success"
		if err != nil {
			msg = "failed"
		}
		notice.HubIns.SendBroadcast(constant.TOPIC_SCRAPER, notice.NoticeResponse{
			Rid:     requestID,
			Event:   constant.EVENT_SCRAPER_SEARCH,
			Message: msg,
			Data: map[string]any{
				"name": s.GetName(),
			},
		})
	}()
	if s == nil {
		return
	}

	param, err := json.Marshal(input)
	if err != nil {
		zaplog.L().Error("刮削参数序列化失败", zap.Any("param", input), zap.Error(err))
		return
	}

	task := &model.Task{
		RequestID:   requestID,
		ScraperName: s.GetName(),
		Type:        model.TaskTypeSearch,
		Param:       string(param),
		Status:      model.TaskStatusWait,
		StartAt:     time.Now(),
	}
	err = data.GetDataFactory().Task().Create(ctx, task, nil)
	if err != nil {
		zaplog.L().Error("刮削失败", zap.Any("param", input), zap.Error(err))
		return
	}

	defer func() {
		err = data.GetDataFactory().Task().Update(ctx, task, nil)
		if err != nil {
			zaplog.L().Error("更新task失败", zap.Uint("task ID", task.ID), zap.Error(err))
			return
		}
	}()

	items, err := s.Search(input.Keyword, input.Page)
	task.EndAt = time.Now()
	task.Duration = int64(task.EndAt.Sub(task.StartAt).Seconds())
	if err != nil {
		task.Status = model.TaskStatusFailed
		zaplog.L().Error("刮削失败", zap.Any("param", input), zap.Error(err))
		return
	}
	data, err := json.Marshal(items)
	if err != nil {
		task.Status = model.TaskStatusFailed
		zaplog.L().Error("刮削失败", zap.Any("param", input), zap.Error(err))
		return
	}
	task.Result = string(data)
	task.Status = model.TaskStatusSucceed
}
