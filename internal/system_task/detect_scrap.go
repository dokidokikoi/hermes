package systemtask

import (
	"context"
	"fmt"
	"izumi/constant"
	"izumi/db"
	"izumi/internal/service"
	"izumi/model"
	"izumi/scraper/event"
	"izumi/scraper/vndb"
	"izumi/utils"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/dokidokikoi/go-common/gopool"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/dokidokikoi/go-common/notice"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// StartAutoDetectScrap 恢复中断的自动探测刮削任务（模板同 StartLoad）：
// 取最新的 running 任务执行，把其余同类型 running 任务标记为 canceled。
func StartAutoDetectScrap() {
	ts, err := db.GetStore().SystemTask().List(context.Background(), &model.SystemTask{
		Type:  model.SystemTaskTypeDetectScrap,
		State: model.SystemTaskStateRunning,
	}, &meta.ListOption{Order: "id desc"})
	if err != nil {
		zaplog.L().Error("system detect scrap task error", zap.Error(err))
		return
	}
	if len(ts) == 0 {
		return
	}
	// 其余 running 的同类型任务视为残留，标记为 canceled
	err = db.GetStore().SystemTask().UpdateByWhere(
		context.Background(),
		&meta.WhereNode{
			Conditions: []*meta.Condition{
				{
					Field:    "id",
					Operator: meta.NOTEQUAL,
					Value:    ts[0].ID,
				},
			},
			Next: &meta.WhereNode{
				Conditions: []*meta.Condition{
					{
						Field:    "state",
						Operator: meta.EQUAL,
						Value:    model.SystemTaskStateRunning,
					},
				},
				Next: &meta.WhereNode{
					Conditions: []*meta.Condition{
						{
							Field:    "type",
							Operator: meta.EQUAL,
							Value:    model.SystemTaskTypeDetectScrap,
						},
					},
				},
			},
		},
		&model.SystemTask{
			State: model.SystemTaskStateCanceled,
		},
		nil,
	)
	if err != nil {
		zaplog.L().Error("system update error", zap.Error(err))
		return
	}

	AutoDetectScrap(ts[0])
}

// AutoDetectScrap 扫描游戏库中未刮削的目录，用目录名在 vndb 搜索，
// 为每个目录派发一个独立的 SystemTaskTypeScrap 子任务（复用 AutoScrap）。
// 父任务在派发完所有子任务后即标记完成，子任务各自独立执行与回报状态。
func AutoDetectScrap(t *model.SystemTask) {
	gopool.CtxGo(context.Background(), func() {
		var e error
		var result string
		defer func() {
			state := model.SystemTaskStateDone
			if e != nil {
				zaplog.L().Error("system detect scrap task error", zap.Error(e))
				state = model.SystemTaskStateFailed
			}
			if err := db.GetStore().SystemTask().Update(context.Background(), &model.SystemTask{
				ID:     t.ID,
				State:  state,
				Result: result,
			}, nil); err != nil {
				zaplog.L().Error("system update error", zap.Error(err))
			}

			if err := notice.HubIns.SendBroadcast("", notice.NoticeResponse{
				Rid:     uuid.NewString(),
				Event:   utils.ConcatEvent(constant.TOPIC_SCRAPER, constant.EVENT_SCRAPER_AUTOSCRAP),
				Success: e == nil,
				Message: result,
			}); err != nil {
				zaplog.L().Error("send notify error", zap.Error(err))
			}
		}()

		// 读取系统策略中的游戏库路径
		p, err := db.GetStore().Policy().Get(context.Background(), &model.Policy{Key: model.SystemPolicy{}.Key()}, nil)
		if err != nil {
			e = err
			result = err.Error()
			return
		}
		sp, err := model.Parse[model.SystemPolicy](p.Policy)
		if err != nil {
			e = err
			result = err.Error()
			return
		}

		// 收集所有游戏库下未刮削（无 info.json）的目录
		library := service.NewSrv(db.GetStore()).Library()
		var pending []service.PathInfo
		for _, lib := range sp.GameLibrary {
			infos, err := library.Ls(context.Background(), lib, true, false)
			if err != nil {
				zaplog.L().Error("get game library error", zap.String("library", lib), zap.Error(err))
				continue
			}
			pending = append(pending, infos...)
		}
		if len(pending) == 0 {
			result = "no games to scrap"
			return
		}

		// vndb 必须已被 PolicyEffect 注册；未注册则无法自动匹配
		vndbScraper, ok := event.GameScraperMap[vndb.Name]
		if !ok {
			e = fmt.Errorf("vndb scraper not registered, check scraper policy")
			result = e.Error()
			return
		}

		// 为每个未刮削目录并发派发子任务：
		//   游戏目录内按后缀找文件名（回退目录名）-> 清洗 -> vndb 搜索取首条
		//   -> 创建 SystemTaskTypeScrap 子任务 -> AutoScrap
		var (
			dispatched int32
			skipped    int32
			wg         sync.WaitGroup
		)
		for _, info := range pending {
			wg.Add(1)
			go func(info service.PathInfo) {
				defer wg.Done()
				dirName := filepath.Base(info.Path)
				rawName := resolveGameName(info, sp.ScrapNameExts)
				// 先从原始名提取版本号（清洗会删除版本号），再清洗名字用于搜索
				version := utils.ExtractVersion(rawName)
				name := utils.CleanGameName(rawName)
				if name == "" {
					// 清洗后可能为空（如整个名字都是括号内容），回退原始目录名
					name = dirName
				}
				items, err := vndbScraper.SearchGame(name, 1)
				if err != nil {
					atomic.AddInt32(&skipped, 1)
					zaplog.L().Error("vndb search failed", zap.String("dir", dirName), zap.String("query", name), zap.Error(err))
					return
				}
				if len(items) == 0 || items[0].URl == "" {
					atomic.AddInt32(&skipped, 1)
					zaplog.L().Warn("vndb search has no match", zap.String("dir", dirName), zap.String("query", name))
					return
				}

				child := &model.SystemTask{
					Type:  model.SystemTaskTypeScrap,
					State: model.SystemTaskStateRunning,
					Param: model.SystemTaskParam{
						ScrapObjs: []model.ScrapObj{{Name: vndb.Name, Url: items[0].URl}},
						Path:      info.Path,
						Version:   version,
					},
				}
				if err := db.GetStore().SystemTask().Create(context.Background(), child, nil); err != nil {
					atomic.AddInt32(&skipped, 1)
					zaplog.L().Error("create child scrap task failed", zap.String("dir", dirName), zap.Error(err))
					return
				}
				atomic.AddInt32(&dispatched, 1)
				AutoScrap(child)
			}(info)
		}
		wg.Wait()

		result = fmt.Sprintf("dispatched %d scrap tasks (%d skipped)", dispatched, skipped)
	})
}

// resolveGameName 从游戏目录内的文件名提取游戏名：
// 按 exts 顺序在目录下找第一个后缀匹配的文件，取其去后缀的文件名；
// 找不到任何匹配文件则回退目录名。name 原样返回（未清洗），由调用方清洗。
func resolveGameName(info service.PathInfo, exts []string) string {
	for _, ext := range exts {
		for _, c := range info.Child {
			if c.IsDir {
				continue
			}
			if strings.EqualFold(filepath.Ext(c.Path), ext) {
				return strings.TrimSuffix(filepath.Base(c.Path), filepath.Ext(c.Path))
			}
		}
	}
	return filepath.Base(info.Path)
}
