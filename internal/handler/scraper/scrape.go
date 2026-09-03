package scraper

import (
	"context"
	"izumi/db"
	"izumi/internal/handler"
	systemtask "izumi/internal/system_task"
	"izumi/model"

	"github.com/dokidokikoi/go-common/middleware"

	"github.com/google/uuid"
)

func (h Handler) Scrape(ctx context.Context, input *handler.ScraperDetailReq, op *middleware.PreHandleOptions) (string, error) {
	if input.RequestID == "" {
		input.RequestID = uuid.New().String()
	}
	requestID := input.RequestID
	_, err := h.srv.Scrap().Scrap(ctx, requestID, input.Objs)

	return requestID, err
}

type AutoScrapeReq struct {
	Objs    []model.ScrapObj `json:"objs"`
	Path    string           `json:"path"`
	Version string           `json:"version"`
	Name    string           `json:"name"`
}

func (h Handler) AutoScrape(ctx context.Context, input *AutoScrapeReq, op *middleware.PreHandleOptions) (string, error) {
	t := &model.SystemTask{
		Type:  model.SystemTaskTypeScrap,
		State: model.SystemTaskStateRunning,
		Param: model.SystemTaskParam{
			ScrapObjs: input.Objs,
			Path:      input.Path,
			Version:   input.Version,
			Name:      input.Name,
		},
	}
	err := db.GetStore().SystemTask().Create(ctx, t, nil)
	if err != nil {
		return "", err
	}
	systemtask.AutoScrap(t)
	return "", nil
}

// AutoDetectScrape 扫描游戏库中未刮削的目录，自动匹配 vndb 并派发刮削子任务。
// 入参为空（路径来自系统策略中的游戏库配置）。
func (h Handler) AutoDetectScrape(ctx context.Context, input *struct{}, op *middleware.PreHandleOptions) (string, error) {
	t := &model.SystemTask{
		Type:  model.SystemTaskTypeDetectScrap,
		State: model.SystemTaskStateRunning,
	}
	if err := db.GetStore().SystemTask().Create(ctx, t, nil); err != nil {
		return "", err
	}
	systemtask.AutoDetectScrap(t)
	return "", nil
}
