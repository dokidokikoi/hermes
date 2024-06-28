package category

import (
	"hermes/db/data"
	"hermes/model"
	"net/http"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h Handler) Del(ctx *gin.Context) {
	ids := []uint{}
	if err := ctx.ShouldBindJSON(&ids); err != nil {
		zaplog.L().Error("参数校验错误", zap.Error(err))
		core.WriteResponse(ctx, errors.ApiErrValidation, nil)
		return
	}

	categoryIDs := []*model.Category{}
	for _, id := range ids {
		categoryIDs = append(categoryIDs, &model.Category{
			ID: id,
		})
	}

	errs := data.GetDataFactory().Category().DeleteCollection(ctx, categoryIDs, nil)
	if len(errs) > 0 {
		zaplog.L().Error("删除失败", zap.Error(errs[0]))
		core.WriteResponse(ctx, &errors.APIError{Code: 1, StatusCode: http.StatusOK, Message: errs[0].Error()}, nil)
		return
	}
	core.WriteResponse(ctx, nil, nil)
}
