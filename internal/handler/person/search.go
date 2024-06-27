package person

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
	var input handler.PersonListReq
	if err := ctx.ShouldBindJSON(&input); err != nil {
		zaplog.L().Error(errors.ApiErrValidation.Message, zap.Error(err))
		core.WriteResponse(ctx, errors.ApiErrValidation, nil)
		return
	}
	var q query.PageQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		core.WriteResponse(ctx, errors.ApiErrValidation, nil)
		return
	}

	total, list, err := h.srv.Person().Search(ctx, input, q.GetListOption(), service.PersonBasicSearchNode...)
	if err != nil {
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}

	vos := make([]handler.StaffVo, len(list))
	for i := range list {
		vos[i] = handler.StaffVo{
			ID:        list[i].ID,
			Name:      list[i].Name,
			Alias:     list[i].Alias,
			Cover:     list[i].Cover,
			Images:    list[i].Images,
			Tags:      list[i].Tags,
			Summary:   list[i].Summary,
			Gender:    list[i].Gender.String(),
			CreatedAt: list[i].CreatedAt,
		}
	}

	core.WriteListResponse(ctx, nil, total, vos)
}
