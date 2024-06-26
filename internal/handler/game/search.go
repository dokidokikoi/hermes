package game

import (
	"hermes/internal/handler"
	"hermes/internal/service"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/dokidokikoi/go-common/query"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h Handler) Search(ctx *gin.Context) {
	var input handler.GameListReq
	if err := ctx.ShouldBindJSON(&input); err != nil {
		zaplog.L().Error(errors.ApiErrValidation.Message, zap.Error(err))
		core.WriteResponse(ctx, errors.ApiErrValidation, nil)
		return
	}
	var q query.PageQuery
	ctx.ShouldBindQuery(&q)

	total, vos, err := h.srv.Game().Search(ctx, input, q.GetListOption(), service.GameBasicSearchNode...)
	if err != nil {
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}

	core.WriteListResponse(ctx, nil, total, vos)
}
