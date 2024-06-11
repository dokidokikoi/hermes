package scraper

import (
	"hermes/db/data"
	"hermes/model"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h Handler) Bind(ctx *gin.Context) {
	var input model.RefGameInstance
	if err := ctx.ShouldBindJSON(&input); err != nil {
		zaplog.L().Error("参数校验错误", zap.Error(err))
		core.WriteResponse(ctx, errors.ApiErrValidation, nil)
		return
	}
	cnt, err := data.GetDataFactory().RefGameInstance().Count(ctx, &model.RefGameInstance{RequestID: input.RequestID}, nil)
	if err != nil {
		zaplog.L().Error("获取刮削关联游戏数据失败", zap.Error(err))
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}
	if cnt > 0 {
		err := data.GetDataFactory().RefGameInstance().Update(ctx, &input, nil)
		if err != nil {
			zaplog.L().Error("更新刮削关联游戏数据失败", zap.Error(err))
			core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
			return
		}
	} else {
		err := data.GetDataFactory().RefGameInstance().Create(ctx, &input, nil)
		if err != nil {
			zaplog.L().Error("插入刮削关联游戏数据失败", zap.Error(err))
			core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
			return
		}
	}
	core.WriteResponse(ctx, nil, nil)
}
