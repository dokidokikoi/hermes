package scraper

import (
	"encoding/json"
	"hermes/db/data"
	"hermes/internal/handler"
	"hermes/model"
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

func (h Handler) Scrape(ctx *gin.Context) {
	var input handler.ScraperDetailReq
	if err := ctx.ShouldBindJSON(&input); err != nil {
		zaplog.L().Error("参数校验错误", zap.Error(err))
		core.WriteResponse(ctx, errors.ApiErrValidation, nil)
		return
	}

	if input.RequestID == "" {
		input.RequestID = uuid.New().String()
	}
	requestID := input.RequestID
	if input.Path != "" {
		err := data.GetDataFactory().RefGameInstance().Create(ctx, &model.RefGameInstance{
			RequestID: requestID,
			Path:      input.Path,
			Version:   input.Version,
		}, nil)
		if err != nil {
			zaplog.L().Error("保存任务关联游戏错误", zap.Error(err))
			core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
			return
		}
	}
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
			task := &model.Task{
				RequestID:   requestID,
				ScraperName: req.Name,
				Type:        model.TaskTypeDetail,
				Param:       req.Url,
				Status:      model.TaskStatusWait,
				StartAt:     time.Now(),
			}
			err := data.GetDataFactory().Task().Create(ctx, task, nil)
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
			task.Status = model.TaskStatusSuccessed
		})
	}
	core.WriteResponse(ctx, nil, requestID)
}
