package scraper

import (
	"errors"
	"hermes/db/data"
	"hermes/model"

	"github.com/dokidokikoi/go-common/core"
	comm_errors "github.com/dokidokikoi/go-common/errors"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (h Handler) GetRef(ctx *gin.Context) {
	requestID := ctx.Query("request_id")
	if requestID == "" {
		zaplog.L().Error("参数缺失")
		core.WriteResponse(ctx, comm_errors.ApiErrValidation, nil)
		return
	}
	res, err := data.GetDataFactory().RefGameInstance().Get(ctx, &model.RefGameInstance{RequestID: requestID}, nil)
	if err != nil {
		zaplog.L().Error("获取关联游戏失败", zap.String("request id", requestID), zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			core.WriteResponse(ctx, comm_errors.ApiErrRecordNotFound, nil)
			return
		}
		core.WriteResponse(ctx, comm_errors.ApiErrSystemErr, nil)
		return
	}
	core.WriteResponse(ctx, nil, res)
}
