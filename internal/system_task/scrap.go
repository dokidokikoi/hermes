package systemtask

import (
	"context"
	"encoding/json"
	"izumi/constant"
	"izumi/db"
	"izumi/internal/handler"
	"izumi/internal/service"
	"izumi/model"
	"izumi/scraper"
	"izumi/scraper/dlsite"
	"izumi/scraper/twodfan"
	"izumi/scraper/vndb"
	"izumi/utils"
	"strings"
	"sync/atomic"
	"time"

	"github.com/abadojack/whatlanggo"
	"github.com/dokidokikoi/go-common/gopool"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/dokidokikoi/go-common/notice"
	"github.com/dokidokikoi/go-common/tools"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func StartAutoScrap() {
	ts, err := db.GetStore().SystemTask().List(context.Background(), &model.SystemTask{
		Type:  model.SystemTaskTypeScrap,
		State: model.SystemTaskStateRunning,
	}, &meta.ListOption{Order: "id desc"})
	if err != nil {
		zaplog.L().Error("system download task error", zap.Error(err))
		return
	}
	if len(ts) == 0 {
		return
	}

	for _, t := range ts {
		AutoScrap(t)
	}
}

func AutoScrap(t *model.SystemTask) {
	gopool.CtxGo(context.Background(), func() {
		var result string
		var e error
		defer func() {
			var state = model.SystemTaskStateDone
			if e != nil {
				zaplog.L().Error("system scrap task error", zap.Error(e))
				state = model.SystemTaskStateFailed
			}
			err := db.GetStore().SystemTask().Update(context.Background(), &model.SystemTask{
				ID:     t.ID,
				State:  state,
				Result: result,
			}, nil)
			if err != nil {
				zaplog.L().Error("system update error", zap.Error(err))
			}

			err = notice.HubIns.SendBroadcast("", notice.NoticeResponse{
				Rid:     uuid.NewString(),
				Event:   utils.ConcatEvent(constant.TOPIC_SCRAPER, constant.EVENT_SCRAPER_AUTOSCRAP),
				Success: e == nil,
				Message: result,
			})
			if err != nil {
				zaplog.L().Error("send notify error", zap.Error(err))
			}
		}()

		flag := false
		for _, obj := range t.Param.ScrapObjs {
			if obj.Name == vndb.Name {
				flag = true
			}
		}
		if !flag {
			result = "auto scrap need vndb item"
			e = errors.New(result)
			return
		}

		var (
			proccess int32
			total    int32 = 10000
			step     int32 = int32(5000 / len(t.Param.ScrapObjs))
		)
		ctxWithCancel, cancel := context.WithCancel(context.TODO())
		defer func() {
			cancel()
			notice.HubIns.SendBroadcast("", notice.NoticeResponse{
				Rid:   uuid.NewString(),
				Event: utils.ConcatEvent(constant.TOPIC_SCRAPER, constant.EVENT_SCRAPER_AUTOSCRAPING),
				Data: map[string]any{
					"task_id":  t.ID,
					"proccess": total,
					"total":    total,
				},
			})
		}()
		go func() {
			ticker := time.NewTicker(time.Millisecond * 500)
			for {
				select {
				case <-ticker.C:
					notice.HubIns.SendBroadcast("", notice.NoticeResponse{
						Rid:   uuid.NewString(),
						Event: utils.ConcatEvent(constant.TOPIC_SCRAPER, constant.EVENT_SCRAPER_AUTOSCRAPING),
						Data: map[string]any{
							"task_id":  t.ID,
							"proccess": atomic.LoadInt32(&proccess),
							"total":    total,
						},
					})
				case <-ctxWithCancel.Done():
					return
				}
			}
		}()
		requestID := uuid.NewString()
		wait, err := service.NewScrap(db.GetStore()).Scrap(context.Background(), requestID, t.Param.ScrapObjs, func(scraperName string, success bool) {
			atomic.AddInt32(&proccess, step)
		})
		if err != nil {
			e = err
			result = err.Error()
			return
		}
		wait.Wait()

		list, err := db.GetStore().Task().List(context.Background(), &model.Task{
			RequestID: requestID,
			Status:    model.TaskStatusSucceed,
		}, nil)
		if err != nil {
			e = err
			result = err.Error()
			return
		}
		itemM := map[string][]scraper.GameItem{}
		for _, l := range list {
			item := scraper.GameItem{}
			err := json.Unmarshal([]byte(l.Result), &item)
			if err != nil {
				zaplog.L().Error("system scrap task error", zap.Error(err))
				continue
			}
			itemM[item.ScraperName] = append(itemM[item.ScraperName], item)
		}

		items, ok := itemM[vndb.Name]
		if !ok || len(items) == 0 {
			result = "vndb item scrap failed"
			return
		}
		delete(itemM, vndb.Name)

		gVo := items[0].GameVo
		aliasSet := tools.NewSet(gVo.Alias...)
		tagM := map[string]struct{}{}
		characterM := map[string]map[int]struct{}{}
		staffM := map[string]map[int]struct{}{}
		for _, t := range gVo.Tags {
			tagM[t.Name] = struct{}{}
		}
		for i, c := range gVo.Characters {
			name := utils.ToLowerNoSpace(c.Name)
			if _, ok := characterM[name]; !ok {
				characterM[name] = map[int]struct{}{}
			}
			characterM[name][i] = struct{}{}
			for _, n := range c.Alias {
				n := utils.ToLowerNoSpace(n)
				if _, ok := characterM[n]; !ok {
					characterM[n] = map[int]struct{}{}
				}
				characterM[n][i] = struct{}{}
			}
		}
		for i, s := range gVo.Staff {
			name := utils.ToLowerNoSpace(s.Name)
			if _, ok := staffM[name]; !ok {
				staffM[name] = map[int]struct{}{}
			}
			staffM[name][i] = struct{}{}
			for _, n := range s.Alias {
				n := utils.ToLowerNoSpace(n)
				if _, ok := staffM[n]; !ok {
					staffM[n] = map[int]struct{}{}
				}
				staffM[n][i] = struct{}{}

			}
		}
		for _, item := range itemM {
			for _, i := range item {
				for _, a := range i.Alias {
					aliasSet.Add(a)
				}
				for _, t := range i.Tags {
					if _, ok := tagM[t.Name]; !ok {
						tagM[t.Name] = struct{}{}
						gVo.Tags = append(gVo.Tags, t)
					}
				}
				gVo.Images = append(gVo.Images, i.Images...)
				gVo.Links = append(gVo.Links, i.Links...)
				if len(gVo.Brands) == 0 {
					gVo.Brands = append(gVo.Brands, i.Brands...)
				}
				if (gVo.Category == nil || gVo.Category.Name == "") && i.Category != nil && i.Category.Name != "" {
					i.Category.Name = strings.ToUpper(i.Category.Name)
					nameBuilder := strings.Builder{}
					for _, r := range i.Category.Name {
						if r >= rune('A') && r <= rune('Z') {
							nameBuilder.WriteRune(r)
						}
					}
					i.Category.Name = nameBuilder.String()
					gVo.Category = i.Category
				}
				if gVo.IssueDate.IsZero() {
					gVo.IssueDate = i.IssueDate
				}
				if gVo.Story == "" {
					gVo.Story = i.Story
				} else if i.Story != "" {
					if needReplace(gVo.Story, i.Story) {
						gVo.Story = i.Story
					}
				}
				if gVo.DlCode == "" {
					gVo.DlCode = i.DlCode
				}
				if gVo.JanCode == "" {
					gVo.JanCode = i.JanCode
				}
				if gVo.Price == "" {
					gVo.Price = i.Price
				}
				for _, c := range i.Characters {
					name := utils.ToLowerNoSpace(c.Name)
					if m, ok := characterM[name]; ok && len(m) == 1 {
						for k := range m {
							gVo.Characters[k].Images = append(gVo.Characters[k].Images, c.Images...)
							if c.Cover != "" {
								gVo.Characters[k].Images = append(gVo.Characters[k].Images, c.Cover)
							}
							if gVo.Characters[k].Gender == model.UnKnown {
								gVo.Characters[k].Gender = c.Gender
							}
							if gVo.Characters[k].Rlation == "" {
								gVo.Characters[k].Rlation = c.Rlation
							}
							if gVo.Characters[k].Summary == "" {
								gVo.Characters[k].Summary = c.Summary
							} else if c.Summary != "" {
								if needReplace(gVo.Characters[k].Summary, c.Summary) {
									gVo.Characters[k].Summary = c.Summary
								}
							}
						}
					} else {
						_, ok := characterM[name]
						if !ok {
							characterM[name] = map[int]struct{}{}
						}
						characterM[name][len(gVo.Characters)] = struct{}{}
						gVo.Characters = append(gVo.Characters, c)
					}
				}
				for _, s := range i.Staff {
					name := utils.ToLowerNoSpace(s.Name)
					if m, ok := staffM[name]; ok && len(m) == 1 {
						for k := range m {
							if gVo.Staff[k].Cover == "" {
								gVo.Staff[k].Cover = s.Cover
							} else if s.Cover != "" {
								gVo.Staff[k].Images = append(gVo.Staff[k].Images, s.Cover)
							}
							gVo.Staff[k].Images = append(gVo.Staff[k].Images, s.Images...)
							gVo.Staff[k].Links = append(gVo.Staff[k].Links, s.Links...)
							if gVo.Staff[k].Gender == model.UnKnown {
								gVo.Staff[k].Gender = s.Gender
							}
							if len(gVo.Staff[k].Relation) == 0 {
								gVo.Staff[k].Relation = s.Relation
							}
							if gVo.Staff[k].Summary != "" {
								gVo.Staff[k].Summary = s.Summary
							} else if s.Summary != "" {
								if needReplace(gVo.Staff[k].Summary, s.Summary) {
									gVo.Staff[k].Summary = s.Summary
								}
							}

						}
					} else {
						_, ok := staffM[name]
						if !ok {
							staffM[name] = map[int]struct{}{}
						}
						staffM[name][len(gVo.Staff)] = struct{}{}
						gVo.Staff = append(gVo.Staff, s)
					}
				}
			}
		}
		gVo.Alias = aliasSet.Slice()

		if vs, ok := itemM[dlsite.Name]; ok {
			for _, v := range vs {
				if v.Story != "" {
					gVo.Story = v.Story
					break
				}
			}
		}
		if vs, ok := itemM[twodfan.Name]; ok {
			for _, v := range vs {
				if v.Story != "" {
					gVo.Story = v.Story
					break
				}
			}
		}

		staffs := []handler.StaffVo{}
		for _, staff := range gVo.Staff {
			for _, r := range staff.Relation {
				if r != model.PRelationUnknown {
					staffs = append(staffs, staff)
					break
				}
			}
		}
		gVo.Staff = staffs
		gVo.Images = tools.NewSet(gVo.Images...).Slice()
		if items, ok := itemM[twodfan.Name]; ok && len(items) > 0 && items[0].Story != "" {
			gVo.Story = items[0].Story
		}

		srv := service.NewGame(db.GetStore())
		err = srv.UpsertFull(context.Background(), &gVo, &model.GameInstance{
			Path:    t.Param.Path,
			Version: t.Param.Version,
		}, func(step int) {
			atomic.AddInt32(&proccess, int32(step))
		})
		if err != nil {
			e = err
			result = err.Error()
			return
		}
		result = gVo.Name

		go func() {
			err = srv.DownloadInfo(context.Background(), gVo.ID, time.Now())
			if err != nil {
				zaplog.L().With(zap.Uint("game id", gVo.ID)).Error("downloadInfo", zap.Error(err))
			}
		}()
	})
}

func needReplace(src, dst string) bool {
	srcLang := whatlanggo.DetectLang(src).Iso6391()
	dstlang := whatlanggo.DetectLang(dst).Iso6391()
	if srcLang == "zh" {
		if dstlang == "zh" {
			return len(src) < len(dst)
		}
		return false
	}
	if dstlang == "zh" {
		return true
	}
	if srcLang == "ja" {
		if dstlang == "ja" {
			return len(src) < len(dst)
		}
		return false
	}
	if dstlang == "ja" {
		return true
	}
	return false
}
