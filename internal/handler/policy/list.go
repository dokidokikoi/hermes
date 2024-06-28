package policy

import (
	"hermes/db/data"
	"hermes/model"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h Handler) List(ctx *gin.Context) {
	ps, err := data.GetDataFactory().Policy().List(ctx, &model.Policy{}, nil)
	if err != nil {
		zaplog.L().Error("获取策略失败", zap.Error(err))
		core.WriteResponse(ctx, errors.ApiErrDatabaseOp, nil)
		return
	}
	res := gin.H{}
	for _, p := range ps {
		res[p.Key] = p.Policy
	}

	core.WriteResponse(ctx, nil, res)
}
