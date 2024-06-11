package scraper

import (
	"encoding/json"
	"fmt"
	"hermes/db/data"
	"hermes/internal/handler"
	"hermes/model"
	"hermes/scraper"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
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

	if len(list) > 0 {
		if list[0].Type == model.TaskTypeDetail {
			var res []scraper.GameItem
			for _, l := range list {
				item := scraper.GameItem{}
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
			return
		} else if list[0].Type == model.TaskTypeSearch {
			res := map[string][]scraper.SearchItem{}
			for _, l := range list {
				items := []scraper.SearchItem{}
				err := json.Unmarshal([]byte(l.Result), &items)
				if err != nil {
					zaplog.L().Error("解析任务失败", zap.String("request id", input.RequestID), zap.Error(err))
					core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
					return
				}

				res[fmt.Sprintf("%s - %s - %d", l.ScraperName, gjson.Get(l.Param, "keyword").String(), gjson.Get(l.Param, "page").Int())] = items
			}
			core.WriteResponse(ctx, nil, res)
			return
		}
	}
	core.WriteResponse(ctx, nil, nil)
}
