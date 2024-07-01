package policy

import (
	"encoding/json"
	"hermes/db/data"
	"hermes/internal/handler"
	"hermes/model"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h Handler) Update(ctx *gin.Context) {
	var input handler.UpdateProxyReq
	if err := ctx.ShouldBindJSON(&input); err != nil {
		zaplog.L().Error("参数解析错误", zap.Error(err))
		core.WriteResponse(ctx, errors.ApiErrValidation, nil)
		return
	}
	var policy any

	switch input.Key {
	case model.SystemPolicy{}.Key():
		t := new(model.SystemPolicy)
		err := json.Unmarshal([]byte(input.Policy), t)
		if err != nil {
			zaplog.L().Error("参数解析错误", zap.Error(err))
			core.WriteResponse(ctx, errors.ApiErrValidation, nil)
			return
		}
		policy = t
		h.srv.Policy().SystemPolicyEffect(ctx, t)
	case model.ScraperPolicy{}.Key():
		t := new(model.ScraperPolicy)
		err := json.Unmarshal([]byte(input.Policy), t)
		if err != nil {
			zaplog.L().Error("参数解析错误", zap.Error(err))
			core.WriteResponse(ctx, errors.ApiErrValidation, nil)
			return
		}
		policy = t
		h.srv.Policy().ScraperPolicyEffect(ctx, t)
	case model.LanguagePolicy{}.Key():
		t := new(model.LanguagePolicy)
		err := json.Unmarshal([]byte(input.Policy), t)
		if err != nil {
			zaplog.L().Error("参数解析错误", zap.Error(err))
			core.WriteResponse(ctx, errors.ApiErrValidation, nil)
			return
		}
		policy = t
	case model.PlatformPolicy{}.Key():
		t := new(model.PlatformPolicy)
		err := json.Unmarshal([]byte(input.Policy), t)
		if err != nil {
			zaplog.L().Error("参数解析错误", zap.Error(err))
			core.WriteResponse(ctx, errors.ApiErrValidation, nil)
			return
		}
		policy = t
	}
	d, err := json.Marshal(policy)
	if err != nil {
		zaplog.L().Error("json编码失败", zap.Error(err))
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}

	err = data.GetDataFactory().Policy().UpdateByWhere(ctx, &meta.WhereNode{
		Conditions: []*meta.Condition{
			{
				Field:    "key",
				Operator: meta.EQUAL,
				Value:    input.Key,
			},
		},
	}, &model.Policy{Policy: string(d)}, nil)
	if err != nil {
		zaplog.L().Error("更新policy数据库失败", zap.Error(err))
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}
	core.WriteResponse(ctx, nil, "")
}
