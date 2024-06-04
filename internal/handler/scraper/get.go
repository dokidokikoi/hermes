package scraper

import (
	"encoding/json"
	"hermes/db/data"
	"hermes/internal/handler"
	"hermes/model"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h Handler) Get(ctx *gin.Context) {
	var input handler.ScraperGetReq
	if err := ctx.ShouldBindQuery(&input); err != nil {
		zaplog.L().Error("参数校验错误", zap.Error(err))
		core.WriteResponse(ctx, errors.ApiErrValidation, nil)
		return
	}

	list, err := data.GetDataFactory().Task().List(ctx, &model.Task{RequestID: input.RequestID, Status: model.TaskStatusSuccessed}, nil)
	if err != nil {
		zaplog.L().Error("获取任务失败", zap.String("request id", input.RequestID), zap.Error(err))
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}
	res := []handler.GameVo{}
	for _, l := range list {
		item := handler.GameVo{}
		err := json.Unmarshal([]byte(l.Result), &item)
		if err != nil {
			zaplog.L().Error("解析任务失败", zap.String("request id", input.RequestID), zap.Error(err))
			core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
			return
		}
		res = append(res, item)
	}
	rg, _ := data.GetDataFactory().RefGameInstance().Get(ctx, &model.RefGameInstance{RequestID: input.RequestID}, nil)
	core.WriteResponse(ctx, nil, gin.H{"game": rg, "list": res})
}
