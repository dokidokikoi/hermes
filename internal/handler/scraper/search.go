package scraper

import (
	"context"
	"encoding/json"
	"hermes/db/data"
	"hermes/internal/handler"
	"hermes/model"
	"hermes/scraper"
	"hermes/scraper/event"
	"time"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	"github.com/dokidokikoi/go-common/gopool"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (h Handler) Search(ctx *gin.Context) {
	var input handler.ScraperSearchReq
	if err := ctx.ShouldBindJSON(&input); err != nil {
		zaplog.L().Error("参数校验错误", zap.Error(err))
		core.WriteResponse(ctx, errors.ApiErrValidation, nil)
		return
	}

	requestID := uuid.New().String()
	if input.Name == "all" {
		for _, s := range event.GameScraperMap {
			s := s
			gopool.CtxGo(ctx, func() {
				DoSearch(ctx, requestID, input, s)
			})
		}
	} else {
		gopool.CtxGo(ctx, func() {
			DoSearch(ctx, requestID, input, event.GameScraperMap[input.Name])
		})
	}
}

func DoSearch(ctx context.Context, requestID string, input handler.ScraperSearchReq, s scraper.IGameScraper) {
	if s == nil {
		return
	}
	task := &model.Task{
		RequestID:   requestID,
		ScraperName: s.GetName(),
		Type:        model.TaskTypeSearch,
		Status:      model.TaskStatusWait,
		StartAt:     time.Now(),
	}
	err := data.GetDataFactory().Task().Create(ctx, task, nil)
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
	task.Status = model.TaskStatusSuccessed
}
